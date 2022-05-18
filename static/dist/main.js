import { Fzf } from "../node_modules/fzf/dist/fzf.es.js";
import { InputSuggest } from "../node_modules/input-suggest/dist/input-suggest.min.js";
const root = new URL(window.location.href);
export const url = root.origin + "/api/v1/data";
export let searchers;
export let allocations;
main();
async function main() {
    searchers = await fetchPath(url + "/allocate/list");
    allocations = await fetchPath(url + "/allocates");
    const textarea = new InputSuggest("textarea");
    InputSuggest("textarea");
    textarea.setSuggestions(["foo", "bar", "biz"]);
}
// fetchの返り値のPromiseを返す
async function fetchPath(url) {
    return await fetch(url)
        .then((response) => {
        return response.json();
    })
        .catch((response) => {
        return Promise.reject(new Error(`{${response.status}: ${response.statusText}`));
    });
}
export function fzfSearch(list, keyword) {
    const fzf = new Fzf(list, {
        selector: (item) => item.body,
    });
    const entries = fzf.find(keyword);
    const ranking = entries.map((entry) => entry.item);
    return ranking;
}
