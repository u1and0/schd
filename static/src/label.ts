import { addListOption, fetchPath } from "./element.js";
import { fzfSearch, type Searcher } from "./fzf.js";

type AddressMap = Map<string, string>;

const root: URL = new URL(window.location.href);
export const url: string = root.origin + "/api/v1/data";
main();

async function main() {
  const addressMap: AddressMap = await fetchPath(url + "/address");
  addListOption(
    document.getElementById("address-list"),
    Object.keys(addressMap),
  );
  const text: HTMLElement = document.getElementById("to-name");
  const inputChange = () => {
    console.log(text.value);
    const addressText = document.getElementById("to-address");
    addressText.value = addressMap[text.value];
  };
  text?.addEventListener("change", inputChange);
}
