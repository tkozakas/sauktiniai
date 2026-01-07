package api

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	persons, _ := h.client.Fetch(region, start, start+limit-1)
	json.NewEncoder(w).Encode(map[string]any{
		"region":  region,
		"start":   start,
		"count":   len(persons),
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

	persons := h.client.Search(region, query, 50000)
	json.NewEncoder(w).Encode(map[string]any{
		"query":   query,
		"region":  region,
		"count":   len(persons),
		"persons": persons,
	})
}
