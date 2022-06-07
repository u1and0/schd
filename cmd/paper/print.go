package paper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/api"
	"github.com/u1and0/schd/cmd/ctrl"
	"github.com/xuri/excelize/v2"
)

// CreatePrintForm : xlsxに転記する
func CreatePrintForm(c *gin.Context) {
	checkBoxLabel := []string{"作業用 ", "外注用 ", "作業引替 ",
		"外注引替 ", "検査用 ", "協議用 ", "承認用 ", "完成図用 ",
		"見積用 ", "参考用 ", "仕様書添付用 ", "要求元控"}
	c.HTML(http.StatusOK, "print_create.tmpl", gin.H{
		"today":         time.Now().Format("2006/01/02"),
		"section":       ctrl.Config.Section,
		"tableRow":      []int{0, 1, 2, 3, 4, 5, 6, 7},
		"checkBoxLabel": checkBoxLabel,
	})
}

// CreatePrint : xlsxに転記するフォームの表示
func CreatePrint(c *gin.Context) {
	// フォームから読み込み
	o := new(api.PrintOrder)
	if err := c.Bind(&o); err != nil {
		msg := err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":        msg,
			"error":      err,
			"PrintOrder": o,
		})
		return
	}
	fmt.Printf("PrintOrder: %#v\n", o)

	// template XLSX
	f, err := excelize.OpenFile("template/template.xlsx")
	defer f.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "error": err})
		return
	}
	cells := Cells{
		"C5": o.Date.Format(LAYOUT),
		"C6": o.Section,
		"C7": o.OrderNo,
		"C8": o.OrderName,
	}
	for i := 0; i < 8; i++ {
		cells[fmt.Sprintf("B%d", i+11)] = o.Drawing.No[i]
		cells[fmt.Sprintf("C%d", i+11)] = o.Drawing.Name[i]
		cells[fmt.Sprintf("D%d", i+11)] = o.Drawing.Quantity[i]
		cells[fmt.Sprintf("E%d", i+11)] = o.Drawing.Deadline[i].Format(LAYOUT)
		cells[fmt.Sprintf("F%d", i+11)] = o.Drawing.Misc[i]
	}
	sheetName := "複写入力画面"
	if err := cells.SetCellValue(f, sheetName); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "error": err})
	}
	downloadFile("P-0-002Q付表1.xlsx", f, c)
}
