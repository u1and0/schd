import { Fzf } from "../node_modules/fzf/dist/fzf.es.js";
const root = new URL(window.location.href);
export const url = root.origin + "/api/v1/data";
export let searchers;
main();
async function main() {
    fetchAddress(url + "/address");
    searchers = await fetchPath(url + "/allocate/list");
}
// fetchの返り値のPromiseを返す
export async function fetchPath(url) {
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
async function fetchAddress(url) {
    try {
        // 住所録.js から住所一覧をselect option に加える;
        const address = await fetchPath(url);
        if (address === null)
            return;
        Object.keys(address).forEach((key) => {
            const elem = document.querySelector("#to-address");
            if (elem !== null)
                elem.append(`<option value=${key}>${key}</option>`);
        });
    }
    catch (error) {
        console.error(`Error occured (${error})`);
    }
}
