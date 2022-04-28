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
    fetchAddress(url.origin + "/api/v1/data/address");
}
// fetchの返り値のPromiseを返す
function fetchLocatePath(url) {
    return __awaiter(this, void 0, void 0, function* () {
        return yield fetch(url)
            .then((response) => {
            // if (!response.ok) {
            // return Promise.reject(new Error(`{${response.status}: ${response.statusText}`));
            // } else{
            return response.json();
        })
            .catch((response) => {
            return Promise.reject(new Error(`{${response.status}: ${response.statusText}`));
        });
    });
}
function fetchAddress(url) {
    return __awaiter(this, void 0, void 0, function* () {
        try {
            // 住所録.js から住所一覧をselect option に加える;
            const address = yield fetchLocatePath(url);
            Object.keys(address).forEach((key) => {
                // const elem = document.getElementById("to-address")
                const elem = document.querySelector('#to-address');
                elem.append(`<option value=${key}>${key}</option>`);
                // $("#to-address").append(`<option value=${key}>${key}</option>`);
            });
            // address.forEach((h) =>{
            //   $("#to-address").append("<option value=" + h.word + "></option>");
            // });
        }
        catch (error) {
            console.error(`Error occured (${error})`);
        }
    });
}
