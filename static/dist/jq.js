import { url, searchers, fzfSearch, allocations } from "./main.js"

function _appendRow(x){
  tr = $(x).closest("tr")
  newtr = tr.clone()
  tr.after(newtr)
}

function _removeRow(x) {
  if ($("#load-table tr").length > 2) {
    $(x).closest("tr").remove()
  } else {
    alert("これ以上行を削除できません。")
  }
}

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
      $("#order").val(el["生産命令番号"])
      $("#section").val(el["型式"])
    });
  })
);
