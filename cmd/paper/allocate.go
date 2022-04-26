package paper

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/api"
	"github.com/xuri/excelize/v2"
)

const (
	ALPATH = "/mnt/2_Common/04_社内標準/_配車要求表_輸送指示書"
)

type (
	// Section 所属課のmap db/課.jsonから読み取る
	Section map[string]string
	// Allocation 配車要求に必要な情報 htmlから入力
	Allocation struct {
		Date     time.Time `json:"要求年月日" form:"allocate-date" time_format:"2006/01/02"`
		Section  string    `json:"型式" form:"section"`
		Type     string    `json:"輸送便の別" form:"type"`
		Car      string    `json:"車種" form:"car"`
		Cartype  string    `json:"台車" form:"car-type"`
		T        int       `json:"t数" form:"t"`
		Function string    `json:"機能" form:"section"`
		Order    string    `json:"生産命令番号" form:"order"`
	}
	// Y 要求票番号と保存されているディレクトリ
	Y struct {
		Base string
		Dir  string
	}
)

// CreateAllocateForm : xlsxに転記するフォームの表示
func CreateAllocateForm(c *gin.Context) {
	section := new(Section)
	if err := api.UnmarshalJSON(section, api.SECTIONFILE); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "allocate_create.tmpl", gin.H{"section": (*section)})
}

// CreateAllocate : xlsx に転記する
func CreateAllocate(c *gin.Context) {
	o := new(Allocation)
	if err := c.Bind(&o); err != nil {
		msg := err.Error()
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": msg, "error": err})
		return
	}
	fmt.Printf("Allocation: %#v\n", o)

	// template XLSX
	f, err := excelize.OpenFile("template/template.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	sheetName := "入力画面"
	// 要求番号
	reqNo, err := getRequestNo(o.Section)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest,
			gin.H{"msg": err.Error(), "error": err})
		return
	}
	f.SetCellValue(sheetName, "F2", reqNo.Base)
	// 要求年月日
	f.SetCellValue(sheetName, "F3", time.Now().Format("2006年1月2日"))
	f.SetCellValue(sheetName, "F4", o.Section)
	// 輸送便の別
	var s string
	switch o.Type {
	case "仕立便":
		s = `☑仕立便　☐常用便
☐混載便　☐宅配便`
	case "常用便":
		s = `☐仕立便　☑常用便
☐混載便　☐宅配便`
	case "混載便":
		s = `☐仕立便　☐常用便
☑混載便　☐宅配便`
	case "宅配便":
		s = `☐仕立便　☐常用便
☐混載便　☑宅配便`
	}
	f.SetCellValue(sheetName, "F5", s)
	// 生産命令番号
	f.SetCellValue(sheetName, "F7", o.Order)
	// 輸送区間
	from, _ := f.GetCellValue(sheetName, "F8")
	to := "適当な宛先"
	f.SetCellValue(sheetName, "F8", from+to)
	// 送り先
	var m api.AddressMap
	if err := api.UnmarshalJSON(m, api.ADDRESSFILE); err != nil {
		c.IndentedJSON(http.StatusBadRequest,
			gin.H{"msg": err.Error(), "error": err})
		return
	}
	a := strings.Join(m[to], "\n")
	f.SetCellValue(sheetName, "F13", a)

	// f.SaveAs(reqNo.Dir + reqNo.Base + ".xlsx")
	downloadFile(sheetName+".xlsx", f, c)
}

func getRequestNo(sec string) (y Y, err error) {
	section := new(Section)
	if err := api.UnmarshalJSON(section, api.SECTIONFILE); err != nil {
		return Y{}, err
	}
	prefix := (*section)[sec]
	surfix := ".xlsx"
	var n int
	err = filepath.Walk(ALPATH,
		func(path string, info os.FileInfo, err error) error {
			base := filepath.Base(path)
			baseBool := strings.HasPrefix(base, prefix)
			extBool := filepath.Ext(path) == surfix
			fmt.Println(path)
			// baseBool かつ extBool => prefixで始まり、xlsxで終わるファイルのみ対象
			if baseBool && extBool {
				// strconv できるかつ
				if _, err := strconv.Atoi(base[5:10]); err == nil {
					y.Base = base // TB00-76045....xlsx
					y.Dir = filepath.Dir(path)
				}
			}
			return nil
		})
	n, err = strconv.Atoi(y.Base[5:10]) // 76045
	if err != nil {
		return
	}
	fmt.Println(n)
	y.Base = prefix + strconv.Itoa(n+1) + surfix
	return
}
