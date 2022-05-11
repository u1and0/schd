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
	// ALPATH : 配車要求票を保存するルートディレクトリ
	ALPATH = "/mnt/2_Common/04_社内標準/_配車要求表_輸送指示書"
	LAYOUT = "2006年1月2日"
	// 別紙記載に分割する行数制限
	MAXLINE = 4
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

// CreateAllocateForm : xlsxに転記するフォームの表示
func CreateAllocateForm(c *gin.Context) {
	section := new(Section)
	if err := api.UnmarshalJSON(section, api.SECTIONFILE); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var hours []int
	for i := 0; i < 24; i++ {
		hours = append(hours, i)
	}
	minutes := []int{0, 15, 30, 45}
	c.HTML(http.StatusOK, "allocate_create.tmpl", gin.H{
		"today":   time.Now().Format("2006/01/02"),
		"section": (*section),
		"hours":   hours,
		"minutes": minutes,
	})
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "error": err})
		return
	}
	sheetName := "入力画面"
	// 要求番号
	reqNo, err := getRequestNo(o.Section)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "error": err})
		return
	}
	f.SetCellValue(sheetName, "F2", reqNo.Base)
	// 要求年月日
	f.SetCellValue(sheetName, "F3", time.Now().Format(LAYOUT))
	f.SetCellValue(sheetName, "F4", o.Section)
	// 輸送便の別
	s := o.Type
	switch s {
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
	m := new(api.AddressMap)
	if err := api.UnmarshalJSON(m, api.ADDRESSFILE); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "error": err})
		return
	}
	a := strings.Join((*m)[to], "\n")
	f.SetCellValue(sheetName, "F13", a)
	// t数 台車 機能
	t := fmt.Sprintf("%st%s(%s)", strconv.Itoa(o.T), o.Cartype, o.Function)
	f.SetCellValue(sheetName, "F6", t)
	// 積込/到着作業月日/時刻
	x := map[string]interface{}{
		"F9":  o.Load.Date.Format(LAYOUT),
		"F10": fmt.Sprintf("%d時%d分", o.Load.Hour, o.Load.Minute),
		"F11": o.Arrive.Date.Format(LAYOUT),
		"F12": fmt.Sprintf("%d時%d分", o.Arrive.Hour, o.Arrive.Minute),
	}
	for cell, value := range x {
		if err = f.SetCellValue(sheetName, cell, value); err != nil {
			fmt.Println(err)
		}
	}
	// 保険要否
	x = map[string]interface{}{
		"F18": `☐要　☑不要`, // 保険
		"F19": "",       // 保険額
	}
	if o.Insulance != "契約済み" {
		x = map[string]interface{}{
			"F18": `☐要　☑不要`,         // 保険
			"F19": o.InsulancePrice, // 保険額
		}
	}
	for cell, value := range x {
		if err = f.SetCellValue(sheetName, cell, value); err != nil {
			fmt.Println(err)
		}
	}
	// 寸法質量
	l := len(o.Package)
	if l < MAXLINE { // 3行までなら配車要求票に記載
		p := o.Size.Compile()
		fmt.Println(p)
		x := map[string]interface{}{
			"F15": o.Size.ToString(), // 重量・長さなど
			"F16": p.ToString(),      // 荷姿(個数)
			"F17": o.Size.Sum(),      // 総個数
		}
		for cell, value := range x {
			if err = f.SetCellValue(sheetName, cell, value); err != nil {
				fmt.Println(err)
			}
		}
	} else { // 4行以上の場合は別紙に記載
		x := map[string]interface{}{
			"F15": "別紙参照",
			"F16": "別紙参照",
			"F17": o.Size.Sum(),
		}
		for cell, value := range x {
			// 重量・長さなど 荷姿(個数) 総個数
			if err = f.SetCellValue(sheetName, cell, value); err != nil {
				fmt.Println(err)
			}
		}
		sheetName := "配車別紙"
		c := 11
		for i := 0; i < l; i++ {
			n := strconv.Itoa(c + i)
			x := map[string]interface{}{
				"F" + n: o.Package[i],
				"G" + n: fmt.Sprintf("%dx%dx%d", o.Width[i], o.Length[i], o.Hight[i]),
				"J" + n: o.Mass[i],
				"K" + n: o.Method[i],
				"L" + n: o.Quantity[i],
			}
			for cell, value := range x {
				if err = f.SetCellValue(sheetName, cell, value); err != nil {
					fmt.Println(err)
				}
			}
		}
	}

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
			// baseBool かつ extBool => prefixで始まり、xlsxで終わるファイルのみ対象
			// strconvする前でないとディレクトリ検知してbaseのスライスできなくてpanic
			// ifが階層化してしまうが後で最適化する
			if baseBool && extBool {
				// strconv できるかつ
				if _, err := strconv.Atoi(base[5:10]); err == nil {
					m, _ := strconv.Atoi(base[5:10]) // 76045
					// 最も大きい数値がある場所を保存ディレクトリとする
					// max int
					if n < m {
						n = m
						y.Dir = filepath.Dir(path)
					}
				}
			}
			return nil
		})
	if err != nil {
		return
	}
	y.Base = prefix + strconv.Itoa(n+1) + surfix
	return
}
