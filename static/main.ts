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
    Object.keys(address).forEach((key: string) =>{
      // const elem = document.getElementById("to-address")
      const elem = document.querySelector('#to-address')
      elem.append(`<option value=${key}>${key}</option>`)
    });
  } catch (error) {
    console.error(`Error occured (${error})`);
  }
}