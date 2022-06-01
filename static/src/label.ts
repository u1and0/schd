import { addListOption, fetchPath } from "./element.js";
import { fzfSearch, type Searcher } from "./fzf.js";

type AddressMap = Map<string, string>;

const root: URL = new URL(window.location.href);
export const url: string = root.origin + "/api/v1/data";
main();

async function main() {
  const addressMap: AddressMap = await fetchPath(url + "/address");
  console.log(addressMap);
  const list = Object.keys(addressMap);
  const e = document.getElementById("address-list");
  if (e === null) return;
  addListOption(e, list);
}
