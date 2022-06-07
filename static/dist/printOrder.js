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
}
$(document).ready($(function () {
    $("#search-form").keyup(function () {
        $("#search-result > option").remove();
        const value = document.getElementById("search-form").value;
        const result = fzfSearchList(printHistoriesList, value);
        console.log(result);
        for (let i = 0; i < result.length; i++) {
            $("#search-result").append($("<option>")
                .html(result[i])
                .val(i));
        }
    });
    $("#search-result").change(function () {
        const i = $("#search-result").val();
        const el = printHistories[i];
        console.log(el);
        $("#section").val(el["要求元"]);
        $("#order-no").val(el["生産命令番号"]);
        $("#order-name").val(el["生産命令名称"]);
        $("input[name='draw-no']").val(el["図番"][0]);
        $("input[name='draw-name']").val(el["図面名称"][0]);
    });
}));
