export async function fetchPath(url) {
    return await fetch(url)
        .then((response) => {
        return response.json();
    })
        .catch((response) => {
        return Promise.reject(new Error(`{${response.status}: ${response.statusText}`));
    });
}
export function addListOption(element, list) {
    [...new Set(list)].sort().map((item) => {
        const option = document.createElement("option");
        option.text = item;
        option.value = item;
        element.appendChild(option);
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
export function checkboxChangeValue() {
    const checkboxes = document.querySelectorAll("input[type='checkbox']");
    checkboxes.forEach((checkbox) => {
        checkbox.addEventListener("change", () => {
            checkbox.value = checkbox.checked ? "true" : "false";
        });
    });
}
