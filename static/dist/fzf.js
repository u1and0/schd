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
      $("#transport").val(el["輸送情報"]["輸送便の別"])
      $("#transport-no").val(el["輸送情報"]["伝票番号"])
      $("#transport-fee").val(el["輸送情報"]["運賃"])
      $("#car").val(el["クラスボディタイプ"])
      $("#to-name").val(el["宛先情報"]["輸送区間"])
      $("textarea#to-address").val(el["宛先情報"]["宛先住所"])
      $("textarea#package-name").val(el["物品情報"]["物品名称"])
      $("textarea#article").val(el["記事"])
      $("textarea#misc").val(el["注意事項"]["その他"])

      $("#piling").val(el["注意事項"]["平積み"])
      checkToggle("#piling")
      $("#fixing").val(el["注意事項"]["固定"])
      checkToggle("#fixing")
      $("#confirm").val(el["注意事項"]["確認"])
      checkToggle("#confirm")
      $("#bill").val(el["注意事項"]["納品書"])
      checkToggle("#bill")
      $("#debt").val(el["注意事項"]["借用書"])
      checkToggle("#debt")
      $("#ride").val(el["注意事項"]["同乗"])
      checkToggle("#ride")
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
