import type { ListResponse, SearchResponse } from "./types";

const API = process.env.NEXT_PUBLIC_API_URL || "";

export const getList = async (region: number, start = 0, limit = 100, year?: string): Promise<ListResponse> => {
  let url = `${API}/api/list?region=${region}&start=${start}&limit=${limit}`;
  if (year) url += `&year=${year}`;
  const res = await fetch(url);
  if (!res.ok) throw new Error("Failed to fetch");
  return res.json();
};

export const search = async (q: string, region = 6): Promise<SearchResponse> => {
  const res = await fetch(`${API}/api/search?q=${encodeURIComponent(q)}&region=${region}`);
  if (!res.ok) throw new Error("Failed to search");
  return res.json();
};

export const getLastUpdated = async (): Promise<string> => {
  const res = await fetch(`${API}/api/updated`);
  if (!res.ok) return "unknown";
  return res.text();
};
