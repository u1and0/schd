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

// CreatePrintForm : xlsxに転記するフォームの表示
func CreatePrintForm(c *gin.Context) {
	c.HTML(http.StatusOK, "print_create.tmpl", gin.H{
		"today":         time.Now().Format("2006/01/02"),
		"section":       ctrl.Config.Section,
		"tableRow":      [8]int{},
		"checkBoxLabel": api.CheckBoxLabel,
	})
}

// CreatePrint : xlsxに転記する
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
	f, err := excelize.OpenFile("template/template_print.xlsx")
	defer f.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "error": err})
		return
	}
	// 必要事項
	cells := Cells{
		"C5": o.Date.Format(LAYOUT),
		"C6": o.Section,
		"C7": o.OrderNo,
		"C8": o.OrderName,
	}
	// 図番　図面名称
	for i := 0; i < 8; i++ {
		j := i + 11
		cells[fmt.Sprintf("B%d", j)] = o.Drawing.No[i]
		cells[fmt.Sprintf("C%d", j)] = o.Drawing.Name[i]
		if o.Drawing.Name[i] != "" {
			cells[fmt.Sprintf("D%d", j)] = o.Drawing.Quantity[i]
			cells[fmt.Sprintf("E%d", j)] = o.Drawing.Deadline[i].Format(LAYOUT)
		}
		cells[fmt.Sprintf("F%d", j)] = o.Drawing.Misc[i]
	}
	// 用途区分及び配布先等
	for i, b := range o.Require {
		j := i + 21
		if b == "true" {
			/* B%dとの比較が必要 Bはどこから読み取る？？ */
			/* templateに渡している配列も動的に生成したい */
			cells[fmt.Sprintf("C%d", j)] = "〇"
		}
	}
	sheetName := "入力画面"
	if err := cells.SetCellValue(f, sheetName); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "error": err})
	}
	downloadFile("複写要求票P-0-002Q付表1.xlsx", f, c)
}
