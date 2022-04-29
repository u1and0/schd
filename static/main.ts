main();

function main() {
  const url = new URL(window.location.href);
  fetchAddress(url.origin + "/api/v1/data/address");
}

// fetchの返り値のPromiseを返す
async function fetchLocatePath(url: string) {
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

async function fetchAddress(url: string) {
  try {
    // 住所録.js から住所一覧をselect option に加える;
    const address = await fetchLocatePath(url);
    if (address === null) return;
    Object.keys(address).forEach((key: string) => {
      const elem = document.querySelector("#to-address");
      if (elem !== null) elem.append(`<option value=${key}>${key}</option>`);
    });
  } catch (error) {
    console.error(`Error occured (${error})`);
  }
}

const tbl = document.getElementById("load-table") as HTMLTableElement;

function _appendRow() {
  if (tbl === null) return;
  const tr = document.createElement("tr");
  const names: Array<string> = [
    "package",
    "width",
    "length",
    "hight",
    "mass",
    "method",
    "quantity",
  ];
  for (const th of names) {
    const td = document.createElement("td");
    const inp = document.createElement("input");
    if (th === "package") {
      inp.setAttribute("name", th);
      inp.setAttribute("list", "package-list");
      inp.setAttribute("size", "10");
      inp.setAttribute("placeholder", "木箱");
    } else if (th === "method") {
      inp.setAttribute("name", th);
      inp.setAttribute("list", "method-list");
      inp.setAttribute("size", "10");
      inp.setAttribute("placeholder", "フォーク");
    } else {
      inp.setAttribute("name", "width");
      inp.setAttribute("class", "small-number");
      inp.setAttribute("placeholder", "0");
    }
    td.append(inp);
    tr.append(td);
  }
  tbl.appendChild(tr);
}

function _removeRow() {
  if (tbl === null) return;
  const l: number = tbl.rows.length;
  tbl.deleteRow(l - 1);
}
