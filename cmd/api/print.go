package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	printHistories = map[string]PrintOrder{}
)

type (
	// PrintOrder : 複写要求 要求事項
	PrintOrder struct {
		Date      time.Time `json:"要求年月日" form:"date" time_format:"2006/01/02"`
		Section   string    `json:"要求元" form:"section"`
		OrderNo   string    `json:"生産命令番号" form:"order-no"`
		OrderName string    `json:"生産命令名称" form:"order-name"`
		Drawing
		Require []bool `json:"必要箇所" form:"require"`
	}
	// Drawing : 図面番号、枚数、期限
	Drawing struct {
		No       [8]string    `json:"図番" form:"draw-no"`
		Name     [8]string    `json:"図面名称" form:"draw-name"`
		Quantity [8]int       `json:"枚数" form:"quantity"`
		Deadline [8]time.Time `json:"要求期限" form:"deadline" time_format:"2006/01/02"`
		Misc     [8]string    `json:"備考" form:"misc"`
	}
)

func init() {
	var filelist []string
	if err := ctrl.UnmarshalJSONfile(&filelist, "db/printlist.json"); err != nil {
		log.Println(err)
	}
	go func() {
		for _, fullpath := range filelist {
			p, err := NewPrintOrder(fullpath)
			if err != nil {
				fmt.Printf("%s: %v", fullpath, err)
				continue
			}
			fmt.Printf("%s: %#v\n", fullpath, p)
			printHistories[p.Concat()] = *p
		}
	}()
}

func NewPrintOrder(fullpath string) (*PrintOrder, error) {
	// .xlsxファイルのみ対象
	// .xlsファイルが混じるとpanic
	if !strings.HasSuffix(fullpath, `.xlsx`) {
		err := fmt.Errorf("error: not xlsx file %s", fullpath)
		return &PrintOrder{}, err
	}
	// ~$ファイルが混じるとpanic
	if strings.Contains(fullpath, `$`) {
		err := fmt.Errorf("error: invalid file name %s", fullpath)
		return &PrintOrder{}, err
	}
	var (
		sheetName = "入力画面"
		p         = new(PrintOrder)
		f, err    = excelize.OpenFile(fullpath)
	)
	if err != nil {
		// f.Close()
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
	for i := 0; i < 8; i++ {
		j := i + 11
		s, _ := f.GetCellValue(sheetName, fmt.Sprintf("B%d", j))
		p.Drawing.No[i] = s
		s, _ = f.GetCellValue(sheetName, fmt.Sprintf("C%d", j))
		p.Drawing.Name[i] = s
		s, _ = f.GetCellValue(sheetName, fmt.Sprintf("D%d", j))
		n, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		p.Drawing.Quantity[i] = n
		s, _ = f.GetCellValue(sheetName, fmt.Sprintf("F%d", j))
		p.Drawing.Misc[i] = s
	}
	// 用途区分及び配布先等
	for i := 21; i < 32; i++ {
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

// Concat : convert search struct
func (p *PrintOrder) Concat() string {
	body := fmt.Sprintf("%v", p)
	body = trimmer(body)
	return body
}
