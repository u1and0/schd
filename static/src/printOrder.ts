import { fetchPath } from "./element.js";
import { fzfSearchList } from "./fzf.js";
import { checkboxesToggle } from "./element.js";

const root = new URL(window.location.href);
const url: string = root.origin + "/api/v1/data";
main();
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
}

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
const drawNo = document.querySelectorAll("input[name='draw-no']");
const drawName = document.querySelectorAll("input[name='draw-name']");
const drawQuant = document.querySelectorAll("input[name='quantity']");
const drawMisc = document.querySelectorAll("input[name='misc']");
outputElem?.addEventListener("change", (e: Event) => {
  const key: string = e.target.value;
  const order: PrintOrder = printHistories[key];
  console.log(order);
  document.getElementById("section").value = order["要求元"];
  document.getElementById("order-no").value = order["生産命令番号"];
  document.getElementById("order-name").value = order["生産命令名称"];
  drawNo.forEach((elem: HTMLElement, i: number) => {
    elem.value = order["図番"][i];
  });
  drawName.forEach((elem: string, i: number) => {
    elem.value = order["図面名称"][i];
  });
  drawQuant.forEach((elem: string, i: number) => {
    elem.value = order["枚数"][i];
  });
  drawMisc.forEach((elem: string, i: number) => {
    elem.value = order["備考"][i];
  });
  checkboxesToggle(order["必要箇所"]);
});

// checkboxを変更するたびに、checkedされている数を枚数に反映
const checkboxes = document.querySelectorAll("input[type='checkbox']");
checkboxes.forEach((checkbox) => {
  checkbox.addEventListener("change", (e: Event) => {
    let checkedCount = 0;
    checkboxes.forEach((checkbox) => {
      if (checkbox.checked) checkedCount++;
    });
    drawQuant.forEach((q, i) => {
      if (drawNo[i].value !== "") {
        q.value = checkedCount;
      }
    });
  });
});
