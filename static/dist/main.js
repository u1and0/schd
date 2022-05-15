"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
main();
function main() {
    const url = new URL(window.location.href);
    const urll = url.origin + "/api/v1/data";
    fetchAddress(urll + "/address");
    fetchAllocate(urll + "/allocate/list");
}
// fetchの返り値のPromiseを返す
function fetchPath(url) {
    return __awaiter(this, void 0, void 0, function* () {
        return yield fetch(url)
            .then((response) => {
            return response.json();
        })
            .catch((response) => {
            return Promise.reject(new Error(`{${response.status}: ${response.statusText}`));
        });
    });
}
function fetchAllocate(url) {
    return __awaiter(this, void 0, void 0, function* () {
        const searchers = yield fetchPath(url);
        const keywords = ["DD", "りんご"];
        for (const searcher of searchers) {
            for (const keyword of keywords) {
                if (searcher["body"].includes(keyword)) {
                    searcher.match += 1;
                }
            }
        }
        console.log(searchers);
    });
}
function fetchAddress(url) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            // 住所録.js から住所一覧をselect option に加える;
            const address = yield fetchPath(url);
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
    });
}
