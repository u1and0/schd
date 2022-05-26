package paper

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ctrl"
)

// CreatePrintForm : xlsxに転記する
func CreatePrintForm(c *gin.Context) {
	c.HTML(http.StatusOK, "print_create.tmpl", gin.H{
		"today":   time.Now().Format("2006/01/02"),
		"section": ctrl.Config.Section,
	})
}

// CreatePrint : xlsxに転記するフォームの表示
func CreatePrint(c *gin.Context) {
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
