import { fetchPath } from "./element.js";
import { fzfSearchList } from "./fzf.js";
import { checkboxChangeValue, checkboxesToggle } from "./element.js";

const root = new URL(window.location.href);
const url: string = root.origin + "/api/v1/data";
main();
checkboxChangeValue();
let printHistories: Map<string, PrintOrder>;

type PrintOrder = {
  "要求元": string;
  "生産命令番号": string;
  "生産命令名称": string;
  "必要箇所": boolean[];
  "図番": string[];
  "図面名称": string[];
  "枚数": number[];
  "備考": string[];
};

async function main() {
  printHistories = await fetchPath(url + "/print");

  const inputElem: HTMLElements = document.getElementById("search-form");
  const outputElem = document.getElementById("search-result");
  // formへの入力があるたびにoption書き換え
  inputElem?.addEventListener("keyup", () => {
    while (outputElem?.firstChild) { // clear option
      outputElem.removeChild(outputElem.firstChild);
    }
    // Map.keys() メソッドがなぜか機能しないのでとりあえずObject.keys()使った。
    const result: string[] = fzfSearchList(
      Object.keys(printHistories),
      inputElem.value,
    );
    result.forEach((key: string) => {
      const option = document.createElement("option");
      option.text = key;
      option.value = key;
      outputElem?.append(option);
    });
  });

  // fzfの結果を出力したoptionを選択するたびに各フォームの書き換え
  outputElem?.addEventListener("change", (e: Event) => {
    const key: string = e.target.value;
    const order: PrintOrder = printHistories[key];
    console.log(order);
    document.getElementById("section").value = order["要求元"];
    document.getElementById("order-no").value = order["生産命令番号"];
    document.getElementById("order-name").value = order["生産命令名称"];
    const drawNo = document.querySelectorAll("input[name='draw-no']");
    drawNo.forEach((elem: HTMLElement, i: number) => {
      elem.value = order["図番"][i];
    });
    const drawName = document.querySelectorAll("input[name='draw-name']");
    drawName.forEach((elem: string, i: number) => {
      elem.value = order["図面名称"][i];
    });
    const drawQuant = document.querySelectorAll("input[name='quantity']");
    drawQuant.forEach((elem: string, i: number) => {
      elem.value = order["枚数"][i];
    });
    const drawMisc = document.querySelectorAll("input[name='misc']");
    drawMisc.forEach((elem: string, i: number) => {
      elem.value = order["備考"][i];
    });
    checkboxesToggle(order["必要箇所"]);
  });
}
