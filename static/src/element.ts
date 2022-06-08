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

export function checkToggle(id: string): void {
  const elem = document.getElementById(id);
  if (elem === null) return;
  if (elem.value === "true") {
    elem.setAttribute("checked", true);
  } else {
    elem.setAttribute("checked", false);
  }
}

// checkboxをクリックしたときにvalueをtrue / false 切り換える
export function checkboxChangeValue() {
  const checkboxes = document.querySelectorAll("input[type='checkbox']");
  checkboxes.forEach((checkbox) => {
    checkbox.addEventListener("change", () => {
      checkbox.value = checkbox.checked ? "true" : "false";
    });
  });
}

// boolListの値がtrueのとき、checkboxをチェックして、valueをtrueにする
// boolListの値がtrue以外のとき、checkboxをチェックをはずして、valueをfalseにする
// eventListenerのコールバック関数内で使う
export function checkboxesToggle(boolList: boolean[]) {
  const checkboxes = document.querySelectorAll("input[type='checkbox']");
  checkboxes.forEach((checkbox: string, idx: number) => {
    checkbox.checked = boolList[idx] ? true : false;
    checkbox.value = boolList[idx] ? "true" : "false";
  });
}
