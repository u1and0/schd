import { fetchPath } from "./element.js";
import { fzfSearchList } from "./fzf.js";
import { checkboxChangeValue, checkboxesToggle } from "./element.js";
const root = new URL(window.location.href);
const url = root.origin + "/api/v1/data";
main();
checkboxChangeValue();
let printHistories;
async function main() {
    printHistories = await fetchPath(url + "/print");
    const inputElem = document.getElementById("search-form");
    const outputElem = document.getElementById("search-result");
    inputElem?.addEventListener("keyup", () => {
        while (outputElem?.firstChild) {
            outputElem.removeChild(outputElem.firstChild);
        }
        const result = fzfSearchList(Object.keys(printHistories), inputElem.value);
        result.forEach((key) => {
            const option = document.createElement("option");
            option.text = key;
            option.value = key;
            outputElem?.append(option);
        });
    });
    outputElem?.addEventListener("change", (e) => {
        const key = e.target.value;
        const order = printHistories[key];
        console.log(order);
        document.getElementById("section").value = order["要求元"];
        document.getElementById("order-no").value = order["生産命令番号"];
        document.getElementById("order-name").value = order["生産命令名称"];
        const drawNo = document.querySelectorAll("input[name='draw-no']");
        drawNo.forEach((elem, i) => {
            elem.value = order["図番"][i];
        });
        const drawName = document.querySelectorAll("input[name='draw-name']");
        drawName.forEach((elem, i) => {
            elem.value = order["図面名称"][i];
        });
        const drawQuant = document.querySelectorAll("input[name='quantity']");
        drawQuant.forEach((elem, i) => {
            elem.value = order["枚数"][i];
        });
        const drawMisc = document.querySelectorAll("input[name='misc']");
        drawMisc.forEach((elem, i) => {
            elem.value = order["備考"][i];
        });
        checkboxesToggle(order["必要箇所"]);
    });
}
