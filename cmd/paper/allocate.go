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
	"github.com/u1and0/schd/cmd/ctrl"
	"github.com/xuri/excelize/v2"
)

const (
	// MAXLINE 別紙記載に分割する行数制限
	MAXLINE = 4
	// LAYOUT : time parse layout
	LAYOUT = "2006年1月2日"
)

// CreateAllocateForm : xlsxに転記するフォームの表示
func CreateAllocateForm(c *gin.Context) {
	var hours []int
	for i := 0; i < 24; i++ {
		hours = append(hours, i)
	}
	minutes := []int{0, 15, 30, 45}
	c.HTML(http.StatusOK, "allocate_create.tmpl", gin.H{
		"today":   time.Now().Format("2006/01/02"),
		"section": ctrl.Config.Section,
		"hours":   hours,
		"minutes": minutes,
	})
}

// CreateAllocate : xlsx に転記する
func CreateAllocate(c *gin.Context) {
	o := new(api.Allocation)
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
	cells := Cells{
		"F2": reqNo.Base,
		"F3": time.Now().Format(LAYOUT), // // 要求年月日
		"F4": o.Section,
	}
	if err := cells.SetCellValue(f, sheetName); err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "error": err})
	}
	// 輸送便の別
	s := o.Transport
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
	// t数 台車 機能 クラス・ボディータイプ
	f.SetCellValue(sheetName, "F6", o.Car)
	// 生産命令番号
	f.SetCellValue(sheetName, "F7", o.Order)
	// 輸送区間
	from, _ := f.GetCellValue(sheetName, "F8")
	f.SetCellValue(sheetName, "F8", from+o.To.Name)
	// 送り先
	// a := strings.Join(ctrl.Config.AddressMap[to], "\n")
	f.SetCellValue(sheetName, "F13", o.To.Address)
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
	if o.Insulance.Need != "契約済み" {
		x = map[string]interface{}{
			"F18": `☐要　☑不要`,          // 保険
			"F19": o.Insulance.Price, // 保険額
		}
	}
	for cell, value := range x {
		if err = f.SetCellValue(sheetName, cell, value); err != nil {
			fmt.Println(err)
		}
	}
	// 寸法質量
	l := len(o.Style)
	if l < MAXLINE { // 3行までなら配車要求票に記載
		p := o.Package.Compile()
		fmt.Println(p)
		x := map[string]interface{}{
			"F15": o.Package.ToString(), // 重量・長さなど
			"F16": p.ToString(),         // 荷姿(個数)
			"F17": o.Package.Sum(),      // 総個数
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
			"F17": o.Package.Sum(),
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
				"F" + n: o.Style[i],
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

	// Refresh 付表１　配車要求票 and 付表２　輸送指示書
	cells = Cells{
		"T1": "" /*要求番号*/, "T2": "" /*要求年月日*/, "F4": "", /*輸送便の別*/
		"F5": "" /*クラス・ボディタイプ*/, "L4": "" /*要求元*/, "F6": "", /*生産命令番号*/
		"F7": "" /*輸送区間*/, "F8": "" /*積込作業日*/, "F9": "", /*到着指定日*/
		"L8": "" /*積込作業時*/, "L9": "" /*到着指定時*/, "F10": "", /*送り先*/
		"F11": "" /*物品名称*/, "F12": "" /*重量長さなど*/, "F13": "", /*荷姿*/
		"F15": "" /*保険*/, "K15": "" /*保険額*/, "O14": "" /*総個数*/, "Q8": "", /*記事*/
	}
	err = cells.SetCellValue(f, "付表１　配車要求票")
	cells = Cells{
		"T1": "" /*要求番号*/, "T2": "" /*要求年月日*/, "Q3": "", /*要求元*/
		"F3": "" /*輸送便の別*/, "G4": "" /*クラス・ボディタイプ*/, "F5": "", /*生産命令番号*/
		"F6": "" /*輸送区間*/, "F7": "", /*積込作業日*/
		"F8": "" /*到着指定日*/, "L7": "", /*積込指定時*/
		"L8": "" /*到着指定時*/, "F9": "", /*送り先*/
		"F10": "" /*重量長さなど*/, "F11": "" /*荷姿*/, "O12": "" /*総個数*/, "Q7": "", /*記事*/
	}
	err = cells.SetCellValue(f, "付表２　輸送指示書")
	f.SaveAs(fmt.Sprintf("%s/%s_%s.xlsx", ctrl.Config.AllocatePath, reqNo.Base, o.Package.Name))
	downloadFile(sheetName+".xlsx", f, c)
}

func getRequestNo(sec string) (y api.Y, err error) {
	prefix := ctrl.Config.Section[sec]
	surfix := ".xlsx"
	var n int
	err = filepath.Walk(ctrl.Config.AllocatePath,
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
