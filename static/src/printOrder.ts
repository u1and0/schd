import { fetchPath } from "./element.js";
import { fzfSearchList } from "./fzf.js";

declare var $: any;

const root = new URL(window.location.href);
const url: string = root.origin + "/api/v1/data";
let printHistoriesList: string[];
let printHistories: PrintHistory[];
main();

type PrintHistory = {
  "要求年月日": string;
  "要求元": string;
  "生産命令番号": string;
  "生産命令名称": string;
  "必要箇所": string;
  "図番": string[];
  "図面名称": string[];
  "枚数": number[];
  "要求期限": string[];
  "備考": string[];
};

async function main() {
  printHistoriesList = await fetchPath(url + "/print/list");
  printHistories = await fetchPath(url + "/print");

  // formへの入力があるたびにoption書き換え
  const inputElem = document.getElementById("search-form");
  const outputElem = document.getElementById("search-result");
  inputElem?.addEventListener("keyup", () => {
    while (outputElem?.firstChild) { // clear option
      outputElem.removeChild(outputElem.firstChild);
    }
    const result: string[] = fzfSearchList(printHistoriesList, inputElem.value);
    result.forEach((line: string, i: number) => {
      const option = document.createElement("option");
      option.text = line;
      option.value = `${i}`;
      outputElem?.append(option);
    });
  });

  // select要素を選択するたびに各フォームへJSONの値を書込み
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

// $(function () {
//   $("#search-result").change(function () {
//     const i = $("#search-result").val();
//     const el = printHistories[i];
//     console.log(el);
//     $("#section").val(el["要求元"]);
//     $("#order-no").val(el["生産命令番号"]);
//     $("#order-name").val(el["生産命令名称"]);
//     $("input[name='draw-no']").val(el["図番"][0]);
//     $("input[name='draw-name']").val(el["図面名称"][0]);
//     // $("#transport-fee").val(el["輸送情報"]["運賃"])
//     // $("#car").val(el["クラスボディタイプ"])
//     // $("#to-name").val(el["宛先情報"]["輸送区間"])
//     // $("textarea#to-address").val(el["宛先情報"]["宛先住所"])
//     // $("textarea#package-name").val(el["物品情報"]["物品名称"])
//     // $("textarea#article").val(el["記事"])
//     // $("textarea#misc").val(el["注意事項"]["その他"])
//   });
// });
