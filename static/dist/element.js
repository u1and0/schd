export async function fetchPath(url) {
    return await fetch(url)
        .then((response) => {
        return response.json();
    })
        .catch((response) => {
        return Promise.reject(new Error(`{${response.status}: ${response.statusText}`));
    });
}
export function addListOption(obj, listid, property) {
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
export function checkboxChengeValue(id) {
    const checkboxes = document.getElementById(id);
    if (checkboxes === null)
        return;
    checkboxes.addEventListener("change", () => {
        checkboxes.value = checkboxes.checked ? "true" : "false";
    });
}
export function checkToggle(id) {
    if ($(id).val() === "true") {
        $(id).prop("checked", true);
    }
    else {
        $(id).prop("checked", false);
    }
}
