import { fzfSearchList, printHistoriesList, printHistories } from "./main.js"

$(document).ready(
  $(function () {
    $("#search-form").keyup(function () {
      $("#search-result > option").remove(); // init option
      const value = document.getElementById("search-form").value;
      const result = fzfSearchList(printHistoriesList, value);
      console.log(result);
      for (let i=0; i<result.length ; i++){
        $("#search-result").append($("<option>")
          .html(result[i])
          .val(i));
      }
    })
    $("#search-result").change(function () {
      const i = $("#search-result").val();
      const el = printHistories[i]
      console.log(el);
      $("#section").val(el["要求元"])
      $("#order-no").val(el["生産命令番号"])
      $("#order-name").val(el["生産命令名称"])
      $("input[name='draw-no']").val(el["図番"][0])
      $("input[name='draw-name']").val(el["図面名称"][0])
      // $("#transport-fee").val(el["輸送情報"]["運賃"])
      // $("#car").val(el["クラスボディタイプ"])
      // $("#to-name").val(el["宛先情報"]["輸送区間"])
      // $("textarea#to-address").val(el["宛先情報"]["宛先住所"])
      // $("textarea#package-name").val(el["物品情報"]["物品名称"])
      // $("textarea#article").val(el["記事"])
      // $("textarea#misc").val(el["注意事項"]["その他"])
    })
  })
)
