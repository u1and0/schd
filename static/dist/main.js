import { Fzf } from "../node_modules/fzf/dist/fzf.es.js";
const root = new URL(window.location.href);
export const url = root.origin + "/api/v1/data";
export let searchers;
export let allocations;
main();
async function main() {
    searchers = await fetchPath(url + "/allocate/list");
    allocations = await fetchPath(url + "/allocates");
    addCarListOption(allocations);
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
function addCarListOption(obj) {
    const select = document.getElementById("car-list");
    if (select === null) {
        return;
    }
    Object.values(obj).map((item) => {
        const option = document.createElement("option");
        const s = item["クラスボディタイプ"];
        console.log(s);
        option.text = s;
        option.value = s;
        select.appendChild(option);
    });
    // for (const item of Object.values(obj)){
    //   console.log(item)
    //   if ( item===null ) { break }
    //   const i:string = item["クラスボディタイプ"]
    //   console.log(i)
    //   option.text = i
    //   option.value = i
    //   select.appendChild(option)
    // }
}
