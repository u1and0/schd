package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ctrl"
	"github.com/xuri/excelize/v2"
)

const (
	// LAYOUT Excel から読み込む日付フォーマット
	LAYOUT = "2006年1月2日"
)

var (
	printHistories = []PrintOrder{}
)

type (
	PrintOrder struct {
		Date      time.Time `json:"要求年月日" form:"date" time_format:"2006年1月2日"`
		Section   string    `json:"要求元" form:"section"`
		OrderNo   string    `json:"生産命令番号" form:"order-no"`
		OrderName string    `json:"清算命令名称" form:"order-name"`
		Drawing
		Require []bool `json:"必要箇所" form:"require"`
	}
	Drawing struct {
		No       []string    `json:"図番" form:"draw-no"`
		Name     []string    `json:"図面名称" form:"draw-name"`
		Quantity []int       `json:"枚数" form:"quantity"`
		Deadline []time.Time `json:"要求期限" form:"deadline" time_format:"2006年1月2日"`
		Misc     []string    `json:"備考" form:"misc"`
	}
)

func init() {
	var filelist []string
	if err := ctrl.UnmarshalJSONfile(&filelist, "db/printlist.json"); err != nil {
		log.Fatal(err)
	}
	go func() {
		for _, fullpath := range filelist {
			// .xlsxファイルのみ対象
			// .xlsファイルが混じるとpanic
			if !strings.HasSuffix(fullpath, `.xlsx`) {
				continue
			}
			// ~$ファイルが混じるとpanic
			if strings.Contains(fullpath, `$`) {
				continue
			}
			p, err := NewPrintOrder(fullpath)
			if err != nil {
				fmt.Printf("%s: %v", fullpath, err)
				continue
			}
			fmt.Printf("%s: %#v\n", fullpath, p)
			printHistories = append(printHistories, *p)
		}
	}()
}

func NewPrintOrder(fullpath string) (*PrintOrder, error) {
	sheetName := "入力画面"
	p := new(PrintOrder)
	f, err := excelize.OpenFile(fullpath)
	if err != nil {
		return p, err
	}
	defer f.Close()
	if f.GetSheetIndex(sheetName) == -1 { //sheetNameがない
		err := fmt.Errorf("error: sheet %v not exist\n", sheetName)
		return p, err
	}
	p.Unmarshal(f, sheetName)
	return p, err
}

func (p *PrintOrder) Unmarshal(f *excelize.File, sheetName string) {
	var err error
	t, _ := f.GetCellValue(sheetName, "C5")
	p.Date, err = time.Parse(LAYOUT, t)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	p.Section, _ = f.GetCellValue(sheetName, "C6")
	p.OrderNo, _ = f.GetCellValue(sheetName, "C7")
	p.OrderName, _ = f.GetCellValue(sheetName, "C8")
	for i := 11; i < 19; i++ {
		s, _ := f.GetCellValue(sheetName, fmt.Sprintf("B%d", i))
		p.Drawing.No = append(p.Drawing.No, s)
		s, _ = f.GetCellValue(sheetName, fmt.Sprintf("C%d", i))
		p.Drawing.Name = append(p.Drawing.Name, s)
	}
}

// FetchPrintHistories : returns printHistories array by parsing Excel files
func FetchPrintHistories(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, printHistories)
}
