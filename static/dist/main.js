import { Fzf } from "../node_modules/fzf/dist/fzf.es.js";
const root = new URL(window.location.href);
export const url = root.origin + "/api/v1/data";
export let searchers;
export let allocations;
main();
async function main() {
    searchers = await fetchPath(url + "/allocate/list");
    allocations = await fetchPath(url + "/allocates");
    addListOption(allocations, "car-list", "クラスボディタイプ");
    const checkBoxIDs = [
        "piling",
        "fixing",
        "confirm",
        "bill",
        "debt",
        "ride",
    ];
    checkBoxIDs.map((id) => {
        checkboxChengeValue(id);
    });
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
function addListOption(obj, listid, property) {
    const select = document.getElementById(listid);
    if (select === null)
        return;
    const carList = [];
    Object.values(obj).map((item) => {
        carList.push(item[property]);
    });
    // Remove duplicate & sort, then append HTML datalist
    [...new Set(carList)].sort().map((item) => {
        const option = document.createElement("option");
        option.text = item;
        option.value = item;
        select.appendChild(option);
    });
}
function checkboxChengeValue(id) {
    const checkboxes = document.getElementById(id);
    checkboxes.addEventListener("change", () => {
        if (checkboxes.checked) {
            checkboxes.value = "true";
        }
        else {
            checkboxes.value = "false";
        }
    });
}
