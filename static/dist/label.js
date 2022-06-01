import { addListOption, fetchPath } from "./element.js";
const root = new URL(window.location.href);
export const url = root.origin + "/api/v1/data";
main();
async function main() {
    const addressMap = await fetchPath(url + "/address");
    addListOption(document.getElementById("address-list"), Object.keys(addressMap));
    const text = document.getElementById("to-name");
    const inputChange = () => {
        console.log(text.value);
        const addressText = document.getElementById("to-address");
        addressText.value = addressMap[text.value];
    };
    text?.addEventListener("change", inputChange);
}
