package karys

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	http  *http.Client
	cache map[int][]Person
	mu    sync.RWMutex
}

func NewClient() *Client {
	c := &Client{
		http:  &http.Client{Timeout: 30 * time.Second},
		cache: make(map[int][]Person),
	}
	c.loadFromDisk()
	return c
}

func (c *Client) loadFromDisk() {
	for region := 1; region <= 6; region++ {
		data, err := os.ReadFile(fmt.Sprintf("data/region_%d.json", region))
		if err != nil {
			continue
		}
		var persons []Person
		if json.Unmarshal(data, &persons) == nil {
			c.cache[region] = persons
		}
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

func (c *Client) FetchAll(region int) []Person {
	c.mu.RLock()
	if cached, ok := c.cache[region]; ok {
		c.mu.RUnlock()
		return cached
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.cache[region]; ok {
		return cached
	}

	var all []Person
	batch := 500
	max := 60000

	var wg sync.WaitGroup
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
			if p, err := c.Fetch(region, start, end); err == nil && len(p) > 0 {
				results <- p
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

	c.cache[region] = all
	return all
}

func (c *Client) Search(region int, query string) []Person {
	all := c.FetchAll(region)
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

func (c *Client) GetCached(region int) []Person {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache[region]
}

func (c *Client) IsCached(region int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.cache[region]
	return ok
}
