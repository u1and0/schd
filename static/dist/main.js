import { fzfSearch } from "./fzf.js";
const root = new URL(window.location.href);
const url = root.origin + "/api/v1/data";
let searchers;
let allocations;
main();
async function main() {
    searchers = await fetchPath(url + "/allocate/list");
    allocations = await fetchPath(url + "/allocates");
    addListOption(allocations, "car-list", "クラスボディタイプ");
    const checkBoxIDs = [
        "piling",
        "fixing",
        "confirm",
        "bill",
        "debt",
        "ride",
    ];
    checkBoxIDs.map((id) => {
        checkboxChengeValue(id);
    });
}
async function fetchPath(url) {
    return await fetch(url)
        .then((response) => {
        return response.json();
    })
        .catch((response) => {
        return Promise.reject(new Error(`{${response.status}: ${response.statusText}`));
    });
}
function addListOption(obj, listid, property) {
    const select = document.getElementById(listid);
    if (select === null)
        return;
    const carList = [];
    Object.values(obj).map((item) => {
        carList.push(item[property]);
    });
    [...new Set(carList)].sort().map((item) => {
        const option = document.createElement("option");
        option.text = item;
        option.value = item;
        select.appendChild(option);
    });
}
function checkboxChengeValue(id) {
    const checkboxes = document.getElementById(id);
    if (checkboxes === null)
        return;
    checkboxes.addEventListener("change", () => {
        checkboxes.value = checkboxes.checked ? "true" : "false";
    });
}
$(function () {
    $("#search-form").keyup(function () {
        $("#search-result > option").remove();
        const value = document.getElementById("search-form").value;
        const result = fzfSearch(searchers, value);
        for (const r of result) {
            $("#search-result").append($("<option>")
                .html(r.body)
                .val(r.id));
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
function checkToggle(id) {
    if ($(id).val() === "true") {
        $(id).prop("checked", true);
    }
    else {
        $(id).prop("checked", false);
    }
}
