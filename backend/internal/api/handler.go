package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"sauktiniai/backend/internal/karys"
)

type Handler struct {
	client *karys.Client
}

func NewHandler(client *karys.Client) *Handler {
	return &Handler{client}
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok"))
}

func (h *Handler) GetLastUpdated(w http.ResponseWriter, _ *http.Request) {
	data, err := os.ReadFile("data/last_updated.txt")
	if err != nil {
		w.Write([]byte("unknown"))
		return
	}
	w.Write([]byte(strings.TrimSpace(string(data))))
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request) {
	region, _ := strconv.Atoi(r.URL.Query().Get("region"))
	if region < 1 || region > 6 {
		region = 6
	}

	start, _ := strconv.Atoi(r.URL.Query().Get("start"))
	if start < 0 {
		start = 0
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	year := r.URL.Query().Get("year")

	var persons []karys.Person
	var total int

	if h.client.IsCached(region) {
		all := h.client.GetCached(region)

		// Filter by year if specified
		if year != "" {
			var filtered []karys.Person
			for _, p := range all {
				if p.Bdate == year {
					filtered = append(filtered, p)
				}
			}
			all = filtered
		}

		total = len(all)
		end := start + limit
		if end > len(all) {
			end = len(all)
		}
		if start < len(all) {
			persons = all[start:end]
		}
	} else {
		persons, _ = h.client.Fetch(region, start, start+limit-1)
		total = len(persons)
	}

	json.NewEncoder(w).Encode(map[string]any{
		"region":  region,
		"start":   start,
		"count":   len(persons),
		"total":   total,
		"persons": persons,
	})
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "q required", http.StatusBadRequest)
		return
	}

	region, _ := strconv.Atoi(r.URL.Query().Get("region"))
	if region < 1 || region > 6 {
		region = 6
	}

	persons := h.client.Search(region, query)
	json.NewEncoder(w).Encode(map[string]any{
		"query":   query,
		"region":  region,
		"count":   len(persons),
		"persons": persons,
	})
}
