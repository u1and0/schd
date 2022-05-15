import { Fzf } from '../node_modules/fzf/dist/fzf.es.js';
main();
function main() {
    /*fzf*/
    const list = [
        "go",
        "fzf",
        "never read",
        "tokyo",
        "Yokohama",
        "Japan",
        "Kyoto",
        "ringo",
    ];
    const fzf = new Fzf(list);
    const entries = fzf.find("go");
    const ranking = entries.map((entry) => entry.item).join(", ");
    console.log(ranking); // => go, ringo
    /*fzf*/
    const url = new URL(window.location.href);
    const urll = url.origin + "/api/v1/data";
    fetchAddress(urll + "/address");
    fetchAllocate(urll + "/allocate/list");
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
async function fetchAllocate(url) {
    const searchers = await fetchPath(url);
    const keywords = ["TB00"];
    for (const searcher of searchers) {
        for (const keyword of keywords) {
            if (searcher["body"].includes(keyword)) {
                searcher.match += 1;
            }
        }
    }
    searchers.sort((i, j) => {
        const keyI = i.match;
        const keyJ = j.match;
        if (keyI < keyJ)
            return 1;
        if (keyI > keyJ)
            return -1;
        return 0;
    });
    const matched = searchers.filter((e) => e.match > 0);
    for (const m of matched) {
        console.log(m.body);
    }
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
