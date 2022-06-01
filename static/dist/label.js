import { addListOption, fetchPath } from "./element.js";
const root = new URL(window.location.href);
export const url = root.origin + "/api/v1/data";
main();
async function main() {
    const addressMap = await fetchPath(url + "/address");
    console.log(addressMap);
    const list = Object.keys(addressMap);
    const e = document.getElementById("address-list");
    if (e === null)
        return;
    addListOption(e, list);
}
