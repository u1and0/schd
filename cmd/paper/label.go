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
	"github.com/u1and0/schd/cmd/api"
	"github.com/xuri/excelize/v2"
)

// LAYOUT for time format
const LAYOUT = "2006/1/2"

type (
	Order struct {
		Title string `json:"件名" form:"title"`
		Order string `json:"オーダー情報" form:"order"`
		// FITTING string `json:"order" form:"order"`
		Name      string    `json:"品名" form:"name"`
		Type      string    `json:"型式" form:"type"`
		Quantity  int       `json:"数量" form:"quantity"`
		WrapDate  time.Time `json:"包装年月日" form:"wrap-date" time_format:"2006/01/02"`
		ToAddress string    `json:"荷受人" form:"to-address"`
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
	// 住所録から選択した宛先の住所を引く
	m := new(api.AddressMap)
	if err := api.UnmarshalJSON(&m, api.ADDRESSFILE); err != nil {
		c.IndentedJSON(http.StatusBadRequest,
			gin.H{"msg": err.Error(), "error": err})
		return
	}

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
		f.SetCellValue(sheetName, cell, o.WrapDate.Format(LAYOUT))
	}
	a := (*m)[o.ToAddress]
	for _, cell := range []string{"H4", "H16"} {
		f.SetCellValue(sheetName, cell, a.String())
	}

	downloadFile("外装ラベル.xlsx", f, c)
}

func downloadFile(fs string, f *excelize.File, c *gin.Context) {
	// Save file
	filepath := "./result.xlsx"
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
		f, _ := ar.Create(fs)
		io.Copy(f, buf)
		return false
	})
}
