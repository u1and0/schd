// fetchの返り値のPromiseを返す
export async function fetchPath(url: string): Promise<any> {
  return await fetch(url)
    .then((response) => {
      return response.json();
    })
    .catch((response) => {
      return Promise.reject(
        new Error(`{${response.status}: ${response.statusText}`),
      );
    });
}

export function addListOption(element: HTMLElement, list: string[]): void {
  // Remove duplicate & sort, then append HTML datalist
  [...new Set(list)].sort().map((item) => {
    const option = document.createElement("option");
    option.text = item;
    option.value = item;
    element.appendChild(option);
  });
}

export function checkboxChengeValue(id: string) {
  const checkboxes: HTMLElement | null = document.getElementById(id);
  if (checkboxes === null) return;
  checkboxes.addEventListener("change", () => {
    // valueはstringの"true","false"
    // Boolean のtrue, falseではない。
    // これはgolangサーバー側でunmarshalするときに"true", "false"という
    // 文字列をいい感じにサーバー側でbool値として解釈してくれるため。
    checkboxes.value = checkboxes.checked ? "true" : "false";
  });
}

export function checkToggle(id: string) {
  if ($(id).val() === "true") {
    $(id).prop("checked", true);
  } else {
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

export function checkboxesToggle(require: boolean[]) {
  const checkboxes = document.querySelectorAll("input[type='checkbox']");
  checkboxes.forEach((checkbox: string, idx: number) => {
    checkbox.value = require[idx];
  });
}
