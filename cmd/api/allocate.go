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
	"github.com/u1and0/schd/cmd/ctrl"
	"github.com/xuri/excelize/v2"
)

var (
	allocations = make(Allocations, 2000)
)

const (
	// LAYOUT : time parse layout
	LAYOUT = "2006年1月2日"
	// ALPATH : 配車要求票を保存するルートディレクトリ
	// ALPATH = "./test"
	ALPATH = "/mnt/2_Common/04_社内標準/_配車要求表_輸送指示書"
	// for win
	// ALPATH = "../../../../../../../../../2_Common/04_社内標準/_配車要求表_輸送指示書"
)

type (
	// Allocations list of Allocation
	Allocations map[ctrl.ID]Allocation
	// Allocation 配車要求に必要な情報 htmlから入力
	Allocation struct {
		Date      time.Time `json:"要求年月日" form:"allocate-date" time_format:"2006/01/02"`
		Section   string    `json:"型式" form:"section"`
		Type      string    `json:"輸送便の別" form:"type"`
		Car       `json:"車両情報"`
		Order     string `json:"生産命令番号" form:"order"`
		To        `json:"宛先情報" form:"to"`
		Load      `json:"積込情報" form:"load"`
		Arrive    `json:"到着情報" form:"arrive"`
		Package   `json:"物品情報" form:"package"`
		Insulance `json:"保険情報" form:"insulance"`
		Article   string `json:"記事" form:"article"`
	}
	// Car : 車両情報
	Car struct {
		Type     string `json:"車種" form:"type"`     // トラック, トレーラ
		Truck    string `json:"台車" form:"truck"`    // 平車、箱車
		T        int    `json:"t数" form:"t"`        // 4t, 10t
		Function string `json:"機能" form:"function"` // エアサス
	}
	// To : 宛先
	To struct {
		Name    string `json:"送り元宛先" form:"to-name"`
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
		Style    []string `json:"荷姿" form:"package"`
		Width    []int    `json:"幅" form:"width"`
		Length   []int    `json:"長さ" form:"length"`
		Hight    []int    `json:"高さ" form:"hight"`
		Mass     []int    `json:"重量" form:"mass"`
		Method   []string `json:"荷下ろし方法" form:"method"`
		Quantity []int    `json:"数量" form:"quantity"`
	}
	// Insulance : 保険要否
	Insulance struct {
		Need  string `json:"保険要否" form:"insulance-bool"`
		Price int    `json:"保険額" form:"insulance-price"`
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
	go func() {
		filepath.Walk(ALPATH,
			func(path string, info os.FileInfo, err error) error {
				var (
					allocation = new(Allocation)
					filename   = filepath.Base(path)
					baseBool   = strings.HasPrefix(filename, "TA00-") ||
						strings.HasPrefix(filename, "TB00-") ||
						strings.HasPrefix(filename, "DD00-")
					extBool = filepath.Ext(path) == ".xlsx"
				)
				// baseBool かつ extBool => prefixで始まり、xlsxで終わるファイルのみ対象
				if baseBool && extBool {
					f, _ := excelize.OpenFile(path)
					// idをファイル名から抽出するか、
					// シートから抽出するかどちらでもよい
					id := ctrl.ID(filename[:10])
					// id, err := f.GetCellValue("入力画面", "F2")
					// if err != nil {
					// 	fmt.Printf("%s: %v\n", filename, err.Error())
					// }

					// Excel セル抽出してAllocate型に充てる
					allocation.Parse(f)
					allocations[id] = *allocation
				}
				return nil
			})
	}()
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

// FetchAllocate : returns allocate object by parsing Excel files or gob
func FetchAllocate(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, allocations)
}

func FetchAllocateID(c *gin.Context) {
	id := ctrl.ID(c.Param("id"))
	c.IndentedJSON(http.StatusOK, allocations[id])
}

// Parse : Parsing Excel file value
func (a *Allocation) Parse(f *excelize.File) {
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
	a.Type, _ = f.GetCellValue(sheetName, "F5")
	// s, _ = f.GetCellValue(sheetName, "F6")
	a.Car.Type, _ = f.GetCellValue(sheetName, "F6")
	// ss := strings.Split(s, `(`) // [ 4t平車, ｴｱｻｽ) ]
	// a.Function = ss[1]
	// st := strings.Split(ss[0], `t`) // [ 4, 平車 ]
	// a.T, err = strconv.Atoi(st[0])
	// if err != nil {
	// fmt.Printf("%v\n", err.Error())
	// }
	// a.Cartype = st[1]
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
	a.Insulance.Need, _ = f.GetCellValue(sheetName, "F18")
	s, _ = f.GetCellValue(sheetName, "F19")
	if s != "" {
		a.Insulance.Price, err = strconv.Atoi(s)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
		}
	}
	a.Article, _ = f.GetCellValue(sheetName, "F20")
}

// Concat : convert search struct
func (as *Allocations) Concat() Searchers {
	var (
		i     int
		s     = make(Searchers, len(*as))
		trims = []string{`{`, `}`, `[`, `]`, "☑", "☐要　", "☐不要 ",
			"☐仕立便　", "☐常用便\n", "☐混載便　", "☐宅配便 ",
			"+0000 ", "UTC ", "00:00:00 ", "0000-", "0001-"}
	)
	for id, val := range *as {
		body := fmt.Sprintf("%s %v", id, val)
		for _, trim := range trims {
			body = strings.ReplaceAll(body, trim, "")
		}
		s[i] = Searcher{
			ID:   id,
			Body: body,
			Date: val.Date,
		}
		i++
	}
	return s
}
