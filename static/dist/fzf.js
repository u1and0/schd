import { Fzf } from "../node_modules/fzf/dist/fzf.es.js";
export function fzfSearch(list, keyword) {
    const fzf = new Fzf(list, {
        selector: (item) => item.body,
    });
    const entries = fzf.find(keyword);
    const ranking = entries.map((entry) => entry.item);
    return ranking;
}
export function fzfSearchList(list, keyword) {
    const fzf = new Fzf(list);
    const entries = fzf.find(keyword);
    const ranking = entries.map((entry) => entry.item);
    return ranking;
}
