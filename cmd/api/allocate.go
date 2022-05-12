package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var (
	allocations = make(Allocations, 2000)
)

const (
	// LAYOUT : time parse layout
	LAYOUT = "2006年1月2日"
	// ALPATH : 配車要求票を保存するルートディレクトリ
	ALPATH = "/mnt/2_Common/04_社内標準/_配車要求表_輸送指示書"
)

type (
	// Allocations list of Allocation
	Allocations map[string]Allocation
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
		Load
		Arrive
		Size
		Insulance      string `json:"保険" form:"insulance"`
		InsulancePrice int    `json:"保険額" form:"insulance-price"`
		Article        string `json:"記事" form:"article"`
	}
	Load struct {
		Date   time.Time `json:"積込作業月日" form:"load-date" time_format:"2006/01/02"`
		Hour   int       `json:"積込指定時" form:"load-hour"`
		Minute int       `json:"積込指定分" form:"load-minute"`
	}
	Arrive struct {
		Date   time.Time `json:"到着作業月日" form:"arrive-date" time_format:"2006/01/02"`
		Hour   int       `json:"到着指定時" form:"arrive-hour"`
		Minute int       `json:"到着指定分" form:"arrive-minute"`
	}
	// Size 荷姿・寸法・重量
	Size struct {
		Package  []string `json:"荷姿" form:"package"`
		Width    []int    `json:"幅" form:"width"`
		Length   []int    `json:"長さ" form:"length"`
		Hight    []int    `json:"高さ" form:"hight"`
		Mass     []int    `json:"重量" form:"mass"`
		Method   []string `json:"荷下ろし方法" form:"method"`
		Quantity []int    `json:"数量" form:"quantity"`
	}
	// Y 要求票番号と保存されているディレクトリ
	Y struct {
		Base string
		Dir  string
	}
	// PackageCount : 荷姿カウンタ
	PackageCount map[string]int
	// Stringfy 表示
	Stringfy interface {
		ToString() string
	}
)

func init() {
	filepath.Walk(ALPATH,
		func(path string, info os.FileInfo, err error) error {
			filename := filepath.Base(path)
			baseBool := strings.HasPrefix(filename, "TA00-") ||
				strings.HasPrefix(filename, "TB00-") ||
				strings.HasPrefix(filename, "DD00-")
			extBool := filepath.Ext(path) == ".xlsx"
			// baseBool かつ extBool => prefixで始まり、xlsxで終わるファイルのみ対象
			if baseBool && extBool {
				f, _ := excelize.OpenFile(path)
				id := filename[:10]
				err = allocations.Parse(f, id)
				fmt.Println(filename, err)
			}
			return nil
		})
}

// Compile : 荷姿によって数量をカウントする
func (s *Size) Compile() PackageCount {
	p := make(PackageCount, len(s.Package))
	for i, k := range s.Package {
		if _, ok := p[k]; !ok {
			p[k] = s.Quantity[i]
		} else {
			p[k] += s.Quantity[i]
		}
	}
	return p
}

// ToString : 表示
func (s *Size) ToString() string {
	var ss []string
	l := len(s.Package)
	for i := 0; i < l; i++ {
		ss = append(ss, fmt.Sprintf("%dx%dx%dmm %dkg", s.Width[i], s.Length[i], s.Hight[i], s.Mass[i]))
	}
	return strings.Join(ss, ", ")
}

// ToString : 表示
func (p *PackageCount) ToString() string {
	var ss []string
	for k, v := range *p {
		ss = append(ss, fmt.Sprintf("%s(%s)", k, strconv.Itoa(v)))
	}
	return strings.Join(ss, ", ")
}

// Sum : 荷姿によらず数量を合計する
func (s *Size) Sum() (n int) {
	for _, q := range s.Quantity {
		n += q
	}
	return
}

// FetchAllocate : returns allocate object by parsing Excel files or gob
func FetchAllocate(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, allocations)
}

func (as *Allocations) Parse(f *excelize.File, id string) error {
	a := new(Allocation)
	sheetName := "入力画面"
	s, err := f.GetCellValue(sheetName, "F3")
	if err != nil {
		return err
	}
	a.Date, err = time.Parse("2006年1月2日", s)
	if err != nil {
		return err
	}
	(*as)[id] = *a
	return err
}
