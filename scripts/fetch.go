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

	date := time.Now().Format("2006-01-02")

	for region := 1; region <= 6; region++ {
		fmt.Printf("Fetching region %d...\n", region)
		persons := fetchAll(region)
		fmt.Printf("Region %d: %d records\n", region, len(persons))

		data, _ := json.Marshal(persons)
		os.WriteFile(fmt.Sprintf("backend/data/region_%d.json", region), data, 0644)
	}

	page := fmt.Sprintf(`"use client";

import { useState, useEffect, useMemo } from "react";
import { getList, search } from "@/lib/api";
import type { Person } from "@/lib/types";

const REGIONS = [
  { id: 1, name: "Alytus" },
  { id: 2, name: "Kaunas" },
  { id: 3, name: "Klaipėda" },
  { id: 4, name: "Panevėžys" },
  { id: 5, name: "Šiauliai" },
  { id: 6, name: "Vilnius" },
];

const currentYear = new Date().getFullYear();
const YEARS = Array.from({ length: 30 }, (_, i) => currentYear - 18 - i);

const LAST_UPDATED = "%s";

export default function Home() {
  const [region, setRegion] = useState(6);
  const [persons, setPersons] = useState<Person[]>([]);
  const [loading, setLoading] = useState(true);
  const [query, setQuery] = useState("");
  const [page, setPage] = useState(0);
  const [isSearch, setIsSearch] = useState(false);
  const [yearFilter, setYearFilter] = useState("");

  const load = async (r: number, p: number) => {
    setLoading(true);
    try {
      const data = await getList(r, p * 100, 100);
      setPersons(data.persons || []);
    } catch {
      setPersons([]);
    }
    setLoading(false);
  };

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) {
      setIsSearch(false);
      load(region, 0);
      setPage(0);
      return;
    }
    setLoading(true);
    setIsSearch(true);
    try {
      const data = await search(query, region);
      setPersons(data.persons || []);
    } catch {
      setPersons([]);
    }
    setLoading(false);
  };

  useEffect(() => {
    if (!isSearch) {
      load(region, page);
    }
  }, [region, page, isSearch]);

  const filtered = useMemo(() => {
    return persons.filter(p => {
      if (yearFilter && p.bdate !== yearFilter) return false;
      return true;
    });
  }, [persons, yearFilter]);

  const reset = () => {
    setPage(0);
    setQuery("");
    setIsSearch(false);
    setYearFilter("");
  };

  return (
    <main className="min-h-screen bg-black text-white">
      <div className="mx-auto max-w-3xl px-4 py-12">
        <h1
          onClick={reset}
          className="cursor-pointer text-center text-4xl font-black tracking-tight hover:opacity-80 transition"
        >
          ŠAUKTINIAI
        </h1>
        <p className="mt-2 text-center text-xs text-zinc-600">
          Atnaujinta: {LAST_UPDATED}
        </p>

        <div className="mt-8 flex flex-wrap justify-center gap-2">
          {REGIONS.map(r => (
            <button
              key={r.id}
              onClick={() => { setRegion(r.id); reset(); }}
              className={`+"`"+`rounded-full px-4 py-2 text-sm font-medium transition ${
                region === r.id ? "bg-white text-black" : "bg-zinc-900 text-zinc-400 hover:text-white"
              }`+"`"+`}
            >
              {r.name}
            </button>
          ))}
        </div>

        <form onSubmit={handleSearch} className="mt-6">
          <input
            type="text"
            value={query}
            onChange={e => setQuery(e.target.value)}
            placeholder="Ieškoti pagal pavardę..."
            className="w-full rounded-xl bg-zinc-900 px-4 py-3 text-white placeholder-zinc-600 outline-none ring-1 ring-zinc-800 focus:ring-zinc-600"
          />
        </form>

        <div className="mt-4 flex gap-3">
          <select
            value={yearFilter}
            onChange={e => setYearFilter(e.target.value)}
            className="flex-1 rounded-lg bg-zinc-900 px-3 py-2 text-sm text-white outline-none ring-1 ring-zinc-800"
          >
            <option value="">Gimimo metai</option>
            {YEARS.map(y => (
              <option key={y} value={y}>{y}</option>
            ))}
          </select>
          {yearFilter && (
            <button
              onClick={() => setYearFilter("")}
              className="rounded-lg bg-zinc-800 px-3 py-2 text-sm text-zinc-400 hover:text-white"
            >
              ✕
            </button>
          )}
        </div>

        {loading ? (
          <div className="mt-12 flex justify-center">
            <div className="h-6 w-6 animate-spin rounded-full border-2 border-zinc-800 border-t-white" />
          </div>
        ) : (
          <>
            <div className="mt-4 text-sm text-zinc-600">
              {filtered.length} įrašų
            </div>

            <div className="mt-2 divide-y divide-zinc-900">
              {filtered.map((p, i) => (
                <div key={`+"`"+`${p.pos}-${i}`+"`"+`} className="flex items-center py-3">
                  <span className="w-12 text-sm text-zinc-600">{p.pos}</span>
                  <span className="font-medium">{p.name} {p.lastname}</span>
                  <span className="ml-3 text-sm text-zinc-600">{p.bdate}</span>
                </div>
              ))}
              {filtered.length === 0 && (
                <div className="py-12 text-center text-zinc-600">Nieko nerasta</div>
              )}
            </div>

            {!isSearch && persons.length > 0 && (
              <div className="mt-6 flex items-center justify-center gap-2">
                <button
                  onClick={() => setPage(p => Math.max(0, p - 1))}
                  disabled={page === 0}
                  className="rounded-lg bg-zinc-900 px-3 py-2 text-sm disabled:opacity-30"
                >
                  ←
                </button>
                <span className="px-4 text-sm text-zinc-500">{page + 1}</span>
                <button
                  onClick={() => setPage(p => p + 1)}
                  disabled={persons.length < 100}
                  className="rounded-lg bg-zinc-900 px-3 py-2 text-sm disabled:opacity-30"
                >
                  →
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </main>
  );
}
`, date)

	os.WriteFile("frontend/src/app/page.tsx", []byte(page), 0644)
	fmt.Printf("Updated LAST_UPDATED to %s\n", date)
	fmt.Println("Done!")
}
