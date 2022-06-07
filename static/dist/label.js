import { addListOption, fetchPath } from "./element.js";
const root = new URL(window.location.href);
export const url = root.origin + "/api/v1/data";
main();
async function main() {
    const addressMap = await fetchPath(url + "/address");
    const addressListElem = document.getElementById("address-list");
    addListOption(addressListElem, Object.keys(addressMap));
    const toname = document.getElementById("to-name");
    toname?.addEventListener("change", () => {
        console.log(toname.value);
        const addressText = document.getElementById("to-address");
        addressText.value = addressMap[toname.value];
    });
}
