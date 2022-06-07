import { fzfSearch, type Searcher } from "./fzf.js";
import {
  addListOption,
  checkboxChengeValue,
  checkToggle,
  fetchPath,
} from "./element.js";

declare var $: any;
const root: URL = new URL(window.location.href);
const url: string = root.origin + "/api/v1/data";
export let searchers: Promise<Searcher[]>;
export let allocations: Promise<unknown>;
export let printHistoriesList: Promise<string[]>;
export let printHistories: Promise<unknown>;
// allocations, printHistoriesが使われないページでも
// ロードされてしまうので、モジュール分割したい
main();

type Allocations = Map<string, Allocation>;
type Allocation = { unknown: unknown };
async function main() {
  searchers = await fetchPath(url + "/allocate/list");
  allocations = await fetchPath(url + "/allocates");
  const list: string[] = [];
  Object.values(allocations).map((item: Allocation) => {
    list.push(item["クラスボディタイプ"]);
  });
  const carElem = document.getElementById("car-list");
  addListOption(carElem, list);
  const checkBoxIDs: Array<string> = [
    "piling",
    "fixing",
    "confirm",
    "bill",
    "debt",
    "ride",
  ];
  checkBoxIDs.map((id: string) => {
    checkboxChengeValue(id);
  });
}

// FZF on keyboard
$(function () {
  $("#search-form").keyup(function () {
    $("#search-result > option").remove(); // reset option
    const value: string = document.getElementById("search-form").value;
    if (value === null) return;
    const result: Searcher[] = fzfSearch(searchers, value);
    for (const r of result) {
      $("#search-result").append(
        $("<option>")
          .html(r.body)
          .val(r.id),
      );
    }
  });
  $("#search-result").change(function () {
    const id = $("#search-result").val();
    const el = allocations[id];
    console.log(el);
    $("#section").val(el["部署"]);
    $("#insulance").val(el["保険額"]);
    $("#transport").val(el["輸送情報"]["輸送便の別"]);
    $("#transport-no").val(el["輸送情報"]["伝票番号"]);
    $("#transport-fee").val(el["輸送情報"]["運賃"]);
    $("#car").val(el["クラスボディタイプ"]);
    $("#to-name").val(el["宛先情報"]["輸送区間"]);
    $("textarea#to-address").val(el["宛先情報"]["宛先住所"]);
    $("textarea#package-name").val(el["物品情報"]["物品名称"]);
    $("textarea#article").val(el["記事"]);
    $("textarea#misc").val(el["注意事項"]["その他"]);

    const checkboxIDProp = {
      "#piling": "平積み",
      "#fixing": "固定",
      "#confirm": "確認",
      "#bill": "納品書",
      "#debt": "借用書",
      "#ride": "同乗",
    };
    Object.entries(checkboxIDProp).forEach(([key, val]) => {
      $(key).val(el["注意事項"][val]);
      checkToggle(key);
    });
  });
});
