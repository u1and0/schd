import { fetchPath } from "./element.js";
import { fzfSearchList } from "./fzf.js";
const root = new URL(window.location.href);
const url = root.origin + "/api/v1/data";
let printHistoriesList;
let printHistories;
main();
async function main() {
    printHistoriesList = await fetchPath(url + "/print/list");
    printHistories = await fetchPath(url + "/print");
    const inputElem = document.getElementById("search-form");
    const outputElem = document.getElementById("search-result");
    inputElem?.addEventListener("keyup", () => {
        while (outputElem?.firstChild) {
            outputElem.removeChild(outputElem.firstChild);
        }
        const result = fzfSearchList(printHistoriesList, inputElem.value);
        result.forEach((line, i) => {
            const option = document.createElement("option");
            option.text = line;
            option.value = `${i}`;
            outputElem?.append(option);
        });
    });
    outputElem?.addEventListener("change", (e) => {
        const idx = e.target.value;
        const val = printHistories[idx];
        console.log(val);
        document.getElementById("section").value = val["要求元"];
        document.getElementById("order-no").value = val["生産命令番号"];
        document.getElementById("order-name").value = val["生産命令名称"];
        document.querySelector("input[name='draw-no']").forEach((elem, i) => {
            elem.value = val["図番"][i];
        });
    });
}
