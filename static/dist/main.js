import { Fzf } from "../node_modules/fzf/dist/fzf.es.js";
main();
async function main() {
    const url = new URL(window.location.href);
    const urll = url.origin + "/api/v1/data";
    fetchAddress(urll + "/address");
    const searchers = await fetchPath(urll + "/allocate/list");
    const keyword = "ZRA-16";
    const result = fzfSearch(searchers, keyword);
    console.log("searchers: ", searchers);
    console.log("result: ", result);
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
function fzfSearch(list, keyword) {
    const fzf = new Fzf(list, {
        selector: (item) => item.body,
    });
    // const input = document.querySelector("#search-form");
    const entries = fzf.find(keyword);
    const ranking = entries.map((entry) => entry.item.body);
    for (const r of ranking) {
        console.log(r);
    }
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
