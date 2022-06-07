package api

import (
	"fmt"
	"log"
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
	addressMap  = make(AddressMap, 10)
)

type (
	// Allocations list of Allocation
	Allocations map[ctrl.ID]Allocation
	// Allocation 配車要求に必要な情報 htmlから入力
	Allocation struct {
		Date      time.Time `json:"要求年月日" form:"allocate-date" time_format:"2006/01/02"`
		Section   string    `json:"部署" form:"section"`
		Transport `json:"輸送情報"`
		Car       string `json:"クラスボディタイプ" form:"car"`
		Order     string `json:"生産命令番号" form:"order"`
		To        `json:"宛先情報" form:"to"`
		Load      `json:"積込情報" form:"load"`
		Arrive    `json:"到着情報" form:"arrive"`
		Package   `json:"物品情報" form:"package"`
		Insulance int    `json:"保険額" form:"insulance"`
		Article   string `json:"記事" form:"article"`
		Check     `json:"注意事項"`
	}
	// Transport : 輸送情報
	Transport struct {
		Name string `json:"輸送便の別" form:"transport"`
		No   string `json:"伝票番号" form:"transport-no"`
		Fee  int    `json:"運賃" form:"transport-fee"`
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
	// Check : 仕立便及び常用便を配車要求する場合の注意事項
	Check struct {
		Piling  bool   `json:"平積み" form:"piling"`
		Fixing  bool   `json:"固定" form:"fixing"`
		Confirm bool   `json:"確認" form:"confirm"`
		Bill    bool   `json:"納品書" form:"bill"`
		Ride    bool   `json:"同乗" form:"ride"`
		Debt    bool   `json:"借用書" form:"debt"`
		Misc    string `json:"その他" form:"misc"`
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
	// AddressMap : excelファイルから読み込んだ過去の住所録
	AddressMap map[string]string
	// Address : 1住所あたり5行まで, 1行あたり15文字まで
	// Excelシートの枠の都合
	// Address []string
)

func init() {
	fullpaths := GetXlsxPath(ctrl.Config.AllocatePath)
	// Sort reverse order
	sort.Sort(sort.Reverse(sort.StringSlice(fullpaths)))
	// Unmarshal xlsx file
	go func() {
		for _, fullpath := range fullpaths {
			// idをファイル名から抽出するか、
			// シートから抽出するかどちらでもよい
			filename := filepath.Base(fullpath)
			id := ctrl.ID(filename[:10])
			a, err := NewAllocation(fullpath)
			if err != nil {
				fmt.Printf("%s: %v", fullpath, err)
				continue
			}
			allocations[id] = *a
			addressMap = setAddressMap(*a)
		}
	}()
}

// GetXlsxPath : return xlsx filepath slice under root
// Append xlsx file fullpath with specific prefix to slice
func GetXlsxPath(root string) (paths []string) {
	err := filepath.Walk(root,
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
	log.Printf("%v\n", err)
	return
}

// NewAllocation : Constructor of Allocation
// Unmarshal by Excel fullpath
func NewAllocation(fullpath string) (*Allocation, error) {
	a := new(Allocation)
	// Excel セル抽出してAllocate型に充てる
	f, err := excelize.OpenFile(fullpath)
	defer f.Close()
	if err == nil {
		a.Unmarshal(f)
	}
	return a, err
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
	a.Transport.Name = trimmer(s)
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
	a.Transport.No, _ = f.GetCellValue(sheetName, "F21")
	if s, _ = f.GetCellValue(sheetName, "F22"); s != "" {
		a.Transport.Fee, err = strconv.Atoi(s)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
		}
	}
	if s, _ := f.GetCellValue(sheetName, "G30"); s == "〇" {
		a.Check.Piling = true
	}
	if s, _ := f.GetCellValue(sheetName, "G31"); s == "〇" {
		a.Check.Fixing = true
	}
	if s, _ := f.GetCellValue(sheetName, "G32"); s == "〇" {
		a.Check.Confirm = true
	}
	if s, _ := f.GetCellValue(sheetName, "G33"); s == "〇" {
		a.Check.Bill = true
	}
	if s, _ := f.GetCellValue(sheetName, "G34"); s == "〇" {
		a.Check.Debt = true
	}
	if s, _ := f.GetCellValue(sheetName, "G35"); s == "〇" {
		a.Check.Ride = true
	}
	a.Check.Misc, _ = f.GetCellValue(sheetName, "B37")
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

// setAddressMap : AllocationsのToフィールドから住所録データを作成する
// api/init() にてallocations[id]へセット後Allocationをセット後、
// ファイル名は新しい順にソートされいる前提で処理が行われるので、
// addressMapに既存の名称があっても住所を上書きしない
func setAddressMap(a Allocation) AddressMap {
	name := a.To.Name
	if _, ok := addressMap[name]; !ok { // 住所上書きしない
		addressMap[name] = a.To.Address
	}
	return addressMap
}
