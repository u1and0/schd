function _appendRow(x){
  if ($(".table-hover tr").length > 20) {
    alert("これ以上行を追加できません。")
    return
  }
  tr = $(x).closest("tr")
  newtr = tr.clone()
  tr.after(newtr)
}

function _removeRow(x) {
  if ($(".table-hover tr").length < 3) {
    alert("これ以上行を削除できません。")
    return
  }
  $(x).closest("tr").remove()
}
