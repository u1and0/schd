import { Fzf } from "../node_modules/fzf/dist/fzf.es.js";

export let searchers: Promise<Searcher[]>;
main();

type Searcher = {
  id: string;
  body: string;
  date: string;
  match: number;
};

async function main() {
  const url: URL = new URL(window.location.href);
  const urll: string = url.origin + "/api/v1/data";
  fetchAddress(urll + "/address");
  searchers = await fetchPath(
    urll + "/allocate/list",
  );
  console.log("searchers: ", searchers);
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
  // const input = document.querySelector("#search-form");
  const entries = fzf.find(keyword);
  const ranking = entries.map((entry: Fzf) => entry.item.body);
  return ranking;
}

async function fetchAddress(url: string) {
  try {
    // 住所録.js から住所一覧をselect option に加える;
    const address = await fetchPath(url);
    if (address === null) return;
    Object.keys(address).forEach((key: string) => {
      const elem = document.querySelector("#to-address");
      if (elem !== null) elem.append(`<option value=${key}>${key}</option>`);
    });
  } catch (error) {
    console.error(`Error occured (${error})`);
  }
}
