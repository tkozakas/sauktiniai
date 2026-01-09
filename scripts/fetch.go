package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Person struct {
	Pos        string `json:"pos"`
	Number     string `json:"number"`
	Name       string `json:"name"`
	Lastname   string `json:"lastname"`
	Bdate      string `json:"bdate"`
	Department string `json:"department"`
	Info       string `json:"info"`
}

func fetch(region, start, end int) ([]Person, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://sauktiniai.karys.lt/list.php?region=%d", region), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:143.0) Gecko/20100101 Firefox/143.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Range-Unit", "items")
	req.Header.Set("Range", fmt.Sprintf("%d-%d", start, end))
	req.Header.Set("Referer", "https://sauktiniai.karys.lt/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var persons []Person
	json.Unmarshal(body, &persons)
	return persons, nil
}

func fetchAll(region int) []Person {
	var all []Person
	batch := 500
	max := 60000

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(chan []Person, max/batch)
	sem := make(chan struct{}, 10)

	for i := 0; i < max/batch; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			start := idx * batch
			end := start + batch - 1
			if p, err := fetch(region, start, end); err == nil && len(p) > 0 {
				mu.Lock()
				results <- p
				mu.Unlock()
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for p := range results {
		all = append(all, p...)
	}

	return all
}

func main() {
	os.MkdirAll("backend/data", 0755)

	for region := 1; region <= 6; region++ {
		fmt.Printf("Fetching region %d...\n", region)
		persons := fetchAll(region)
		fmt.Printf("Region %d: %d records\n", region, len(persons))

		if persons == nil {
			persons = []Person{}
		}
		data, _ := json.Marshal(persons)
		os.WriteFile(fmt.Sprintf("backend/data/region_%d.json", region), data, 0644)
	}

	date := time.Now().Format("2006-01-02")
	os.WriteFile("backend/data/last_updated.txt", []byte(date), 0644)
	fmt.Printf("Last updated: %s\n", date)
	fmt.Println("Done!")
}
