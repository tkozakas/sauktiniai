export interface Person {
  pos: string;
  number: string;
  name: string;
  lastname: string;
  bdate: string;
  department: string;
  info: string;
}

export interface Region {
  id: number;
  name: string;
}

export interface ListResponse {
  region: number;
  start: number;
  count: number;
  persons: Person[];
}

export interface SearchResponse {
  query: string;
  region: number;
  count: number;
  persons: Person[];
}
