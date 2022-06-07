import { addListOption, fetchPath } from "./element.js";
import { fzfSearch, type Searcher } from "./fzf.js";

type AddressMap = Map<string, string>;

const root: URL = new URL(window.location.href);
export const url: string = root.origin + "/api/v1/data";
main();

async function main() {
  const addressMap: AddressMap = await fetchPath(url + "/address");
  const addressListElem = document.getElementById("address-list")
  addListOption(addressListElem, Object.keys(addressMap));
  const toname: HTMLElement = document.getElementById("to-name");
  toname?.addEventListener("change", () => {
    console.log(toname.value);
    const addressText = document.getElementById("to-address");
    addressText.value = addressMap[toname.value];
  });
}
