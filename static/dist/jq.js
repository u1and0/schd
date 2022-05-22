function _appendRow(x){
  tr = $(x).closest("tr")
  newtr = tr.clone()
  tr.after(newtr)
}

function _removeRow(x) {
  if ($("#load-table tr").length > 2) {
    $(x).closest("tr").remove()
  } else {
    alert("これ以上行を削除できません。")
  }
}
