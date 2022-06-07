package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ctrl"
	"github.com/xuri/excelize/v2"
)

const (
	// LAYOUT Excel から読み込む日付フォーマット
	LAYOUT = "2006年1月2日"
)

var (
	printHistories = PrintHistories{}
)

type (
	PrintHistories []PrintOrder
	PrintOrder     struct {
		Section   string `json:"要求元" form:"section"`
		OrderNo   string `json:"生産命令番号" form:"order-no"`
		OrderName string `json:"生産命令名称" form:"order-name"`
		Drawing
		Require []bool `json:"必要箇所" form:"require"`
	}
	Drawing struct {
		No       []string `json:"図番" form:"draw-no"`
		Name     []string `json:"図面名称" form:"draw-name"`
		Quantity []int    `json:"枚数" form:"quantity"`
		// Deadline []time.Time `json:"要求期限" form:"deadline" time_format:"2006年1月2日"`
		Misc []string `json:"備考" form:"misc"`
	}
)

func init() {
	var filelist []string
	if err := ctrl.UnmarshalJSONfile(&filelist, "db/printlist.json"); err != nil {
		log.Println(err)
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
	p.Section, _ = f.GetCellValue(sheetName, "C6")
	p.OrderNo, _ = f.GetCellValue(sheetName, "C7")
	p.OrderName, _ = f.GetCellValue(sheetName, "C8")
	// 図番 図面名称 枚数　要求期限　備考
	for i := 11; i < 19; i++ {
		s, _ := f.GetCellValue(sheetName, fmt.Sprintf("B%d", i))
		p.Drawing.No = append(p.Drawing.No, s)
		s, _ = f.GetCellValue(sheetName, fmt.Sprintf("C%d", i))
		p.Drawing.Name = append(p.Drawing.Name, s)
		s, _ = f.GetCellValue(sheetName, fmt.Sprintf("D%d", i))
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		p.Drawing.Quantity = append(p.Drawing.Quantity, n)
		s, _ = f.GetCellValue(sheetName, fmt.Sprintf("F%d", i))
		p.Drawing.Misc = append(p.Drawing.Misc, s)
	}
	// 用途区分及び配布先等
	for i := 20; i < 32; i++ {
		s, _ := f.GetCellValue(sheetName, fmt.Sprintf("C%d", i))
		if s != "" {
			p.Require = append(p.Require, true)
		} else {
			p.Require = append(p.Require, false)
		}
	}
}

// FetchPrintHistories : returns printHistories object by parsing Excel files
func FetchPrintHistories(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, printHistories)
}

// FetchPrintList : returns printHistories array by parsing Excel files
func FetchPrintList(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, printHistories.Concat())
}

// Concat : convert search struct
func (p *PrintHistories) Concat() []string {
	s := make([]string, len(*p))
	for i, val := range *p {
		body := fmt.Sprintf("%v", val)
		s[i] = trimmer(body)
	}
	return s
}
