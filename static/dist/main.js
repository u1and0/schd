define(["require", "exports", "../node_modules/fzf/dist/fzf.es.js"], function (require, exports, fzf_es_js_1) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    exports.fzfSearch = exports.fetchPath = exports.callAllocation = exports.searchers = exports.url = void 0;
    const root = new URL(window.location.href);
    exports.url = root.origin + "/api/v1/data";
    let allocations;
    main();
    function callAllocation() {
        alert(allocations);
    }
    exports.callAllocation = callAllocation;
    async function main() {
        fetchAddress(exports.url + "/address");
        exports.searchers = await fetchPath(exports.url + "/allocate/list");
        allocations = await fetchPath(exports.url + "/allocates");
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
    exports.fetchPath = fetchPath;
    function fzfSearch(list, keyword) {
        const fzf = new fzf_es_js_1.Fzf(list, {
            selector: (item) => item.body,
        });
        const entries = fzf.find(keyword);
        const ranking = entries.map((entry) => entry.item);
        return ranking;
    }
    exports.fzfSearch = fzfSearch;
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
});
