main()

function main(){
  const url = new URL(window.location.href);
  fetchAddress(url.origin + "/api/v1/data/address");
}

// fetchの返り値のPromiseを返す
function fetchLocatePath(url){
  return fetch(url)
    .then(response =>{
      // if (!response.ok) {
        // return Promise.reject(new Error(`{${response.status}: ${response.statusText}`));
      // } else{
      return response.json();
    })
  .catch(response =>{
    return Promise.reject(
      new Error(`{${response.status}: ${response.statusText}`));
  });
}


async function fetchAddress(url){
  try{
    // 住所録.js から住所一覧をselect option に加える;
    const address = await fetchLocatePath(url);
    Object.keys(address).forEach((key) =>{
      $("#to-address").append(`<option value=${key}>${key}</option>`);
    });
    // address.forEach((h) =>{
    //   $("#to-address").append("<option value=" + h.word + "></option>");
    // });
  } catch(error) {
    console.error(`Error occured (${error})`);
  }
}


