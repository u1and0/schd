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

// CreateForm : xlsxに転記するフォームの表示
func CreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "label_create.tmpl", "")
}

// Create : POST メソッドによりフォームからの情報をxlsxに転記する
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
	sheetName := "外装表示(単品用)"
	for _, cell := range []string{"D3", "D15"} {
		f.SetCellValue(sheetName, cell, o.Title)
	}
	for _, cell := range []string{"D4", "D16"} {
		f.SetCellValue(sheetName, cell, o.Order)
	}
	for _, cell := range []string{"D6", "D18"} {
		f.SetCellValue(sheetName, cell, o.Name)
	}
	for _, cell := range []string{"D7", "D19"} {
		f.SetCellValue(sheetName, cell, o.Name)
	}
	num := strconv.Itoa(o.Quantity) + "SE"
	for _, cell := range []string{"D8", "D20"} {
		f.SetCellValue(sheetName, cell, num)
	}
	for _, cell := range []string{"D9", "D21"} {
		f.SetCellValue(sheetName, cell, o.WrapDate.Format("2006/1/2"))
	}

	// Save file
	filepath := "result/result.xlsx"
	if err := f.SaveAs(filepath); err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(filepath)

	// Download file
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
}
