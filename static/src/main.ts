import { Fzf } from "../node_modules/fzf/dist/fzf.es.js";

const root: URL = new URL(window.location.href);
export const url: string = root.origin + "/api/v1/data";
export let searchers: Promise<Searcher[]>;
export let allocations;
main();

type Searcher = {
  id: string;
  body: string;
  date: string;
  match: number;
};

async function main() {
  searchers = await fetchPath(url + "/allocate/list");
  allocations = await fetchPath(url + "/allocates");
  addCarListOption(allocations);
}

// fetchの返り値のPromiseを返す
async function fetchPath(url: string): Promise<any> {
  return await fetch(url)
    .then((response) => {
      return response.json();
    })
    .catch((response) => {
      return Promise.reject(
        new Error(`{${response.status}: ${response.statusText}`),
      );
    });
}

export function fzfSearch(list: Searcher[], keyword: string): string[] {
  const fzf = new Fzf(list, {
    selector: (item: Searcher) => item.body,
  });
  const entries = fzf.find(keyword);
  const ranking: string[] = entries.map((entry: Fzf) => entry.item);
  return ranking;
}

function addCarListOption(obj): void {
  const select: HTMLElement | null = document.getElementById("car-list");
  if (select === null) return;
  const carList: Array<string> = [];
  Object.values(obj).map((item: unknown) => {
    carList.push(item["クラスボディタイプ"]);
  });
  // Remove duplicate & sort, then append HTML datalist
  [...new Set(carList)].sort().map((item) => {
    const option = document.createElement("option");
    option.text = item;
    option.value = item;
    select.appendChild(option);
  });
}
