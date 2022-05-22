package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ctrl"
	"github.com/xuri/excelize/v2"
)

var (
	allocations = make(Allocations, 2000)
)

type (
	// Allocations list of Allocation
	Allocations map[ctrl.ID]Allocation
	// Allocation 配車要求に必要な情報 htmlから入力
	Allocation struct {
		Date      time.Time `json:"要求年月日" form:"allocate-date" time_format:"2006/01/02"`
		Section   string    `json:"部署" form:"section"`
		Transport string    `json:"輸送便の別" form:"transport"`
		Car       string    `json:"クラスボディタイプ" form:"car"`
		Order     string    `json:"生産命令番号" form:"order"`
		To        `json:"宛先情報" form:"to"`
		Load      `json:"積込情報" form:"load"`
		Arrive    `json:"到着情報" form:"arrive"`
		Package   `json:"物品情報" form:"package"`
		Insulance int    `json:"保険額" form:"insulance"`
		Article   string `json:"記事" form:"article"`
	}
	// To : 宛先
	To struct {
		Name    string `json:"輸送区間" form:"to-name"`
		Address string `json:"宛先住所" form:"to-address"`
	}
	// Load : 積み込み
	Load struct {
		Date   time.Time `json:"積込作業月日" form:"load-date" time_format:"2006/01/02"`
		Hour   int       `json:"積込指定時" form:"load-hour"`
		Minute int       `json:"積込指定分" form:"load-minute"`
		Time   time.Time `json:"積込指定時刻" form:"load-time"`
	}
	// Arrive : 到着
	Arrive struct {
		Date   time.Time `json:"到着作業月日" form:"arrive-date" time_format:"2006/01/02"`
		Hour   int       `json:"到着指定時" form:"arrive-hour"`
		Minute int       `json:"到着指定分" form:"arrive-minute"`
		Time   time.Time `json:"到着指定時刻" form:"arrive-time"`
	}
	// Package 荷姿・寸法・重量
	Package struct {
		Name     string   `json:"物品名称" form:"package-name"`
		Style    []string `json:"荷姿" form:"style"`
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
	// Append xlsx file fullpath with specific prefix to slice
	paths := []string{}
	filepath.Walk(ctrl.Config.AllocatePath,
		func(path string, info os.FileInfo, err error) error {
			var (
				filename = filepath.Base(path)
				baseBool = strings.HasPrefix(filename, "TA00-") ||
					strings.HasPrefix(filename, "TB00-") ||
					strings.HasPrefix(filename, "DD00-")
				extBool = filepath.Ext(path) == ".xlsx"
			)
			// baseBool かつ extBool => prefixで始まり、xlsxで終わるファイルのみ対象
			if baseBool && extBool {
				paths = append(paths, path)
			}
			return nil
		})
	// Sort reverse order
	sort.Sort(sort.Reverse(sort.StringSlice(paths)))
	// Unmarshal xlsx file
	go func(fullpaths []string) {
		for _, fullpath := range fullpaths {
			allocation := new(Allocation)
			// idをファイル名から抽出するか、
			// シートから抽出するかどちらでもよい
			filename := filepath.Base(fullpath)
			id := ctrl.ID(filename[:10])
			// Excel セル抽出してAllocate型に充てる
			f, err := excelize.OpenFile(fullpath)
			defer f.Close()
			if err == nil {
				allocation.Unmarshal(f)
				allocations[id] = *allocation
			}
		}
	}(paths)
}

// Compile : 荷姿によって数量をカウントする
func (p *Package) Compile() PackageCount {
	n := make(PackageCount, len(p.Style))
	for i, k := range p.Style {
		if _, ok := n[k]; !ok {
			n[k] = p.Quantity[i]
		} else {
			n[k] += p.Quantity[i]
		}
	}
	return n
}

// ToString : 表示
func (p *Package) ToString() string {
	var ss []string
	l := len(p.Style)
	for i := 0; i < l; i++ {
		ss = append(ss, fmt.Sprintf("%dx%dx%dmm %dkg", p.Width[i], p.Length[i], p.Hight[i], p.Mass[i]))
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
func (p *Package) Sum() (n int) {
	for _, q := range p.Quantity {
		n += q
	}
	return
}

// FetchAllocates : returns Allocates object by parsing Excel files
func FetchAllocates(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, allocations)
}

// FetchAllocate : returns Allocate object by parsing Excel files by id
func FetchAllocate(c *gin.Context) {
	id := ctrl.ID(c.Param("id")) // Cast
	if allocation, ok := allocations[id]; ok {
		c.IndentedJSON(http.StatusOK, allocation)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v not found", id)})
}

// Unmarshal : Parsing Excel file value
func (a *Allocation) Unmarshal(f *excelize.File) {
	var (
		err       error
		sheetName = "入力画面"
	)
	s, _ := f.GetCellValue(sheetName, "F3")
	a.Date, err = time.Parse("2006年1月2日", s)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	a.Section, _ = f.GetCellValue(sheetName, "F4")
	s, _ = f.GetCellValue(sheetName, "F5")
	a.Transport = trimmer(s)
	a.Car, _ = f.GetCellValue(sheetName, "F6")
	a.Order, _ = f.GetCellValue(sheetName, "F7")
	a.To.Name, _ = f.GetCellValue(sheetName, "F8")
	s, _ = f.GetCellValue(sheetName, "F9")
	a.Load.Date, err = time.Parse("1月2日", s)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	s, _ = f.GetCellValue(sheetName, "F10")
	a.Load.Time, err = time.Parse("15時04分", s)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	s, _ = f.GetCellValue(sheetName, "F11")
	a.Arrive.Date, err = time.Parse("1月2日", s)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	s, _ = f.GetCellValue(sheetName, "F12")
	a.Arrive.Time, err = time.Parse("15時04分", s)
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
	a.To.Address, _ = f.GetCellValue(sheetName, "F13")
	a.Package.Name, _ = f.GetCellValue(sheetName, "F14")
	if s, _ = f.GetCellValue(sheetName, "F19"); s != "" {
		a.Insulance, err = strconv.Atoi(s)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
		}
	}
	a.Article, _ = f.GetCellValue(sheetName, "F20")
}

func trimmer(s string) string {
	trims := []string{`{`, `}`, `[`, `]`, "☑", "☐要　", "☐不要 ",
		"☐仕立便　", "☐常用便\n", "☐混載便　", "　☐宅配便", "☐宅配便 ", "\n☐宅配便",
		"+0000 ", "UTC ", "00:00:00 ", "0000-", "0001-"}
	for _, trim := range trims {
		s = strings.ReplaceAll(s, trim, "")
	}
	return s
}

// Concat : convert search struct
func (as *Allocations) Concat() Searchers {
	var (
		i int
		s = make(Searchers, len(*as))
	)
	for id, val := range *as {
		body := fmt.Sprintf("%s %v", id, val)
		body = trimmer(body)
		s[i] = Searcher{
			ID:   id,
			Body: body,
			Date: val.Date,
		}
		i++
	}
	return s
}
