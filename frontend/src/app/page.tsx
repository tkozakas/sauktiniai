"use client";

import { useState, useEffect } from "react";
import { getList, search, getLastUpdated } from "@/lib/api";
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

export default function Home() {
  const [region, setRegion] = useState(6);
  const [persons, setPersons] = useState<Person[]>([]);
  const [loading, setLoading] = useState(true);
  const [query, setQuery] = useState("");
  const [page, setPage] = useState(0);
  const [total, setTotal] = useState(0);
  const [isSearch, setIsSearch] = useState(false);
  const [yearFilter, setYearFilter] = useState("");
  const [lastUpdated, setLastUpdated] = useState("");

  useEffect(() => {
    getLastUpdated().then(setLastUpdated);
  }, []);

  const load = async (r: number, p: number, year?: string) => {
    setLoading(true);
    try {
      const data = await getList(r, p * 100, 100, year || undefined);
      setPersons(data.persons || []);
      setTotal(data.total || 0);
    } catch {
      setPersons([]);
      setTotal(0);
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
      load(region, page, yearFilter);
    }
  }, [region, page, isSearch, yearFilter]);

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
        {lastUpdated && (
          <p className="mt-2 text-center text-xs text-zinc-600">
            Atnaujinta: {lastUpdated}
          </p>
        )}

        <div className="mt-8 flex flex-wrap justify-center gap-2">
          {REGIONS.map(r => (
            <button
              key={r.id}
              onClick={() => { setRegion(r.id); reset(); }}
              className={`rounded-full px-4 py-2 text-sm font-medium transition ${
                region === r.id ? "bg-white text-black" : "bg-zinc-900 text-zinc-400 hover:text-white"
              }`}
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
            onChange={e => { setYearFilter(e.target.value); setPage(0); }}
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
              {persons.length} iš {total} įrašų
            </div>

            <div className="mt-2 divide-y divide-zinc-900">
              {persons.map((p, i) => (
                <div key={`${p.pos}-${i}`} className="grid grid-cols-[3rem_1fr_4rem] gap-2 items-start py-3">
                  <span className="text-sm text-zinc-600">{p.pos}</span>
                  <div className="min-w-0">
                    <span className="font-medium">{p.name} {p.lastname}</span>
                    {p.info && <p className="text-xs text-zinc-500 truncate">{p.info}</p>}
                  </div>
                  <span className="text-sm text-zinc-600 text-right">{p.bdate}</span>
                </div>
              ))}
              {persons.length === 0 && (
                <div className="py-12 text-center text-zinc-600">Nieko nerasta</div>
              )}
            </div>

            {!isSearch && total > 100 && (
              <div className="mt-6 flex items-center justify-center gap-2">
                <button
                  onClick={() => setPage(p => Math.max(0, p - 1))}
                  disabled={page === 0}
                  className="rounded-lg bg-zinc-900 px-3 py-2 text-sm disabled:opacity-30"
                >
                  ←
                </button>
                <span className="px-4 text-sm text-zinc-500">{page + 1} / {Math.ceil(total / 100)}</span>
                <button
                  onClick={() => setPage(p => p + 1)}
                  disabled={(page + 1) * 100 >= total}
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
