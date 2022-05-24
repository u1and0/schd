import { url, searchers, fzfSearch, allocations } from "./main.js"

// FZF on keyboard
$(document).ready(
  $(function () {
    $("#search-form").keyup(function () {
      $("#search-result > option").remove(); // init option
      const value = document.getElementById("search-form").value;
      const result = fzfSearch(searchers, value);
      for (const r of result){
        $("#search-result").append($("<option>")
          .html(r.body)
          .val(r.id));
      }
    });
    $("#search-result").change(function () {
      const id = $("#search-result").val();
      const el = allocations[id]
      console.log(el);
      $("#section").val(el["部署"])
      $("#insulance").val(el["保険額"])
      $("#transport").val(el["輸送情報"]["輸送便の別"])
      $("#transport-no").val(el["輸送情報"]["伝票番号"])
      $("#transport-fee").val(el["輸送情報"]["運賃"])
      $("#car").val(el["クラスボディタイプ"])
      $("#to-name").val(el["宛先情報"]["輸送区間"])
      $("textarea#to-address").val(el["宛先情報"]["宛先住所"])
      $("textarea#package-name").val(el["物品情報"]["物品名称"])
      $("textarea#article").val(el["記事"])
      $("textarea#misc").val(el["注意事項"]["その他"])

      const checkboxIDProp = {
        "#piling": "平積み",
        "#fixing": "固定",
        "#confirm": "確認",
        "#bill": "納品書",
        "#debt": "借用書",
        "#ride": "同乗",
      }
      Object.entries(checkboxIDProp).forEach(([key, val]) => {
        $(key).val(el["注意事項"][val]);
        checkToggle(key);
      });
    });
  })
)

function checkToggle(id){
  if ($(id).val() === "true") {
    $(id).prop("checked", true);
  } else {
    $(id).prop("checked", false);
  }
}
