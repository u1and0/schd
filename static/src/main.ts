main();

type Searcher = {
  id: string;
  body: string;
  date: string;
  match: number;
};

function main() {
  const url: URL = new URL(window.location.href);
  const urll: string = url.origin + "/api/v1/data";
  fetchAddress(urll + "/address");
  fetchAllocate(urll + "/allocate/list");}

// fetchの返り値のPromiseを返す
async function fetchPath(url: string): Promise<any>{
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

async function fetchAllocate(url:string){
  const searchers = await fetchPath(url);
  const keywords = ["DD", "りんご"]
  for (const searcher of searchers){
    for (const keyword of keywords){
      if (searcher["body"].includes(keyword)){
        searcher.match+=1
      }
    }
  }
  searchers.sort((i: Searcher,j: Searcher) => {
    const keyI = i.match
    const keyJ = j.match
    if (keyI<keyJ) return 1;
    if (keyI>keyJ) return -1;
    return 0
  })
  console.log(searchers);
}

async function fetchAddress(url: string) {
  try {
    // 住所録.js から住所一覧をselect option に加える;
    const address = await fetchPath(url);
    if (address === null) return;
    Object.keys(address).forEach((key: string) => {
      const elem = document.querySelector("#to-address");
      if (elem !== null) elem.append(`<option value=${key}>${key}</option>`);
    });
  } catch (error) {
    console.error(`Error occured (${error})`);
  }
}
