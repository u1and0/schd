package paper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreatePrint: xlsxに転記する
func CreatePrint(c *gin.Context) {
	c.HTML(http.StatusOK, "copy_create.tmpl", gin.H{})
}

// CreatePrintForm : xlsxに転記するフォームの表示
func CreatePrintForm(c *gin.Context) {
	// // フォームから読み込み
	// o := new(api.PrintOrder)
	// if err := c.Bind(&o); err != nil {
	// 	msg := err.Error()
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"msg":        msg,
	// 		"error":      err,
	// 		"Allocation": o,
	// 	})
	// 	return
	// }
	// fmt.Printf("Allocation: %#v\n", o)
	//
	// // template XLSX
	// f, err := excelize.OpenFile("template/template.xlsx")
	// defer f.Close()
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "error": err})
	// 	return
	// }
	// sheetName := "入力画面"
	// downloadFile("P-0-002Q付表1.xlsx", f, c)
}
