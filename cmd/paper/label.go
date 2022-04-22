package paper

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type (
	Order struct {
		Title string `json:"件名" form:"title"`
		Order string `json:"オーダー情報" form:"order"`
		// FITTING string `json:"order" form:"order"`
		Name     string    `json:"品名" form:"name"`
		Type     string    `json:"型式" form:"type"`
		Quantity int       `json:"数量" form:"quantity"`
		WrapDate time.Time `json:"包装年月日" form:"wrap-date" time_format:"2006/01/02"`
	}
	Cell struct {
		要求番号  string
		要求年月日 time.Time
	}
)

func CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "label_create.tmpl", "")
}

func Create(c *gin.Context) {
	// Parse query
	o := new(Order)
	if err := c.Bind(&o); err != nil {
		msg := err.Error()
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": msg, "error": err})
		return
	}
	fmt.Printf("%#v", o)

	// template XLSX
	f, err := excelize.OpenFile("template/template.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// o := Order{"ZRA-16", "13333", 1, time.Date(2020, 1, 9, 0, 0, 0, 0, time.Local)}
	sheetName := "外装表示(単品用)"
	f.SetCellValue(sheetName, "D3", o.Title)
	f.SetCellValue(sheetName, "D15", o.Title)
	f.SetCellValue(sheetName, "D4", o.Order)
	f.SetCellValue(sheetName, "D16", o.Order)
	f.SetCellValue(sheetName, "D6", o.Name)
	f.SetCellValue(sheetName, "D18", o.Name)
	f.SetCellValue(sheetName, "D7", o.Type)
	f.SetCellValue(sheetName, "D19", o.Type)
	num := strconv.Itoa(o.Quantity) + "SE"
	f.SetCellValue(sheetName, "D8", num)
	f.SetCellValue(sheetName, "D20", num)
	f.SetCellValue(sheetName, "D9", o.WrapDate.Format("2006/1/2"))
	f.SetCellValue(sheetName, "D21", o.WrapDate.Format("2006/1/2"))

	// Save file
	filepath := "result/result.xlsx"
	if err := f.SaveAs(filepath); err != nil {
		fmt.Println(err)
		return
	}

	// Download file
	// c.FileAttachment(filepath, "外装ラベル(単体用).xlsx")
	buf, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Stream(func(w io.Writer) bool {
		ar := zip.NewWriter(w)
		defer ar.Close()
		c.Writer.Header().Set("Content-Disposition", "attachmnt; filename=download.zip")
		f, _ := ar.Create("外装ラベル.xlsx")
		io.Copy(f, buf)
		return false
	})

	// HTML
	// msg := fmt.Sprintf("ラベルを作成しました")
	// c.HTML(http.StatusOK, "get.tmpl", gin.H{"msg": msg, "Order": o})
}
