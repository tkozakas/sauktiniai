package karys

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	return &Client{http: &http.Client{Timeout: 30 * time.Second}}
}

func (c *Client) Regions() []Region {
	return []Region{
		{1, "Alytaus RKC"},
		{2, "Kauno RKC"},
		{3, "Klaipėdos RKC"},
		{4, "Panevėžio RKC"},
		{5, "Šiaulių RKC"},
		{6, "Vilniaus RKC"},
	}
}

func (c *Client) Fetch(region, start, end int) ([]Person, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://sauktiniai.karys.lt/list.php?region=%d", region), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:143.0) Gecko/20100101 Firefox/143.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Range-Unit", "items")
	req.Header.Set("Range", fmt.Sprintf("%d-%d", start, end))
	req.Header.Set("Referer", "https://sauktiniai.karys.lt/")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var persons []Person
	json.Unmarshal(body, &persons)
	return persons, nil
}

func (c *Client) FetchAll(region, max int) []Person {
	batch := 500
	batches := (max + batch - 1) / batch

	var wg sync.WaitGroup
	results := make(chan []Person, batches)
	sem := make(chan struct{}, 10)

	for i := 0; i < batches; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			start := idx * batch
			end := start + batch - 1
			if end >= max {
				end = max - 1
			}
			if p, err := c.Fetch(region, start, end); err == nil {
				results <- p
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var all []Person
	for p := range results {
		all = append(all, p...)
	}
	return all
}

func (c *Client) Search(region int, query string, max int) []Person {
	all := c.FetchAll(region, max)
	query = strings.ToLower(query)

	var matched []Person
	for _, p := range all {
		if strings.Contains(strings.ToLower(p.Name), query) ||
			strings.Contains(strings.ToLower(p.Lastname), query) ||
			strings.Contains(p.Number, query) ||
			strings.Contains(p.Bdate, query) {
			matched = append(matched, p)
		}
	}
	return matched
}
