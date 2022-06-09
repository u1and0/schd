import { Fzf } from "../node_modules/fzf/dist/fzf.es.js";

export type Searcher = {
  id: string;
  body: string;
  date: string;
  match: number;
};

export function fzfSearch(list: Searcher[], keyword: string): Searcher[] {
  const fzf = new Fzf(list, {
    selector: (item: Searcher) => item.body,
  });
  const entries = fzf.find(keyword);
  const ranking: Searcher[] = entries.map((entry: Fzf) => entry.item);
  return ranking;
}

export function fzfSearchList(list: string[], keyword: string): string[] {
  const fzf = new Fzf(list);
  const entries = fzf.find(keyword);
  const ranking: string[] = entries.map((entry: Fzf) => entry.item);
  return ranking;
}

