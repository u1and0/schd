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
      $("#transport").val(el["輸送便の別"])
      $("#car").val(el["クラスボディタイプ"])
      $("#to-name").val(el["宛先情報"]["輸送区間"])
      $("textarea#to-address").val(el["宛先情報"]["宛先住所"])
      $("textarea#package-name").val(el["物品情報"]["物品名称"])
      $("textarea#article").val(el["記事"])
    });
  })
)
