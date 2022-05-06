package paper

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/api"
	"github.com/xuri/excelize/v2"
)

type (
	// Order 要求情報
	Order struct {
		Title string `json:"件名" form:"title"`
		Order string `json:"オーダー情報" form:"order"`
		// FITTING string `json:"order" form:"order"`
		Name      string    `json:"品名" form:"name"`
		Type      string    `json:"型式" form:"type"`
		Quantity  int       `json:"数量" form:"quantity"`
		WrapDate  time.Time `json:"包装年月日" form:"wrap-date" time_format:"2006/01/02"`
		ToAddress string    `json:"荷受人" form:"to-address"`
	}
	Cell struct {
		要求番号  string
		要求年月日 time.Time
	}
)

// CreateForm : xlsxに転記するフォームの表示
func CreateForm(c *gin.Context) {
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
	c.HTML(http.StatusOK, "label_create.tmpl", gin.H{
		"today":   time.Now().Format(LAYOUT),
		"section": (*section),
		"hours":   hours,
		"minutes": minutes,
	})
}

// Create : POST メソッドによりフォームからの情報をxlsxに転記する
func Create(c *gin.Context) {
	// template XLSX
	f, err := excelize.OpenFile("template/template.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	// Parse query
	o := new(Order)
	if err := c.Bind(&o); err != nil {
		msg := err.Error()
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": msg, "error": err})
		return
	}
	// 住所録から選択した宛先の住所を引く
	m := new(api.AddressMap)
	if err := api.UnmarshalJSON(&m, api.ADDRESSFILE); err != nil {
		c.IndentedJSON(http.StatusBadRequest,
			gin.H{"msg": err.Error(), "error": err})
		return
	}

	sheetName := "外装表示(単品用)"
	for _, cell := range []string{"D3", "D15"} {
		f.SetCellValue(sheetName, cell, o.Title)
	}
	for _, cell := range []string{"D4", "D16"} {
		f.SetCellValue(sheetName, cell, o.Order)
	}
	for _, cell := range []string{"D6", "D18"} {
		f.SetCellValue(sheetName, cell, o.Name)
	}
	for _, cell := range []string{"D7", "D19"} {
		f.SetCellValue(sheetName, cell, o.Name)
	}
	num := strconv.Itoa(o.Quantity) + "SE"
	for _, cell := range []string{"D8", "D20"} {
		f.SetCellValue(sheetName, cell, num)
	}
	for _, cell := range []string{"D9", "D21"} {
		f.SetCellValue(sheetName, cell, o.WrapDate.Format("2006/1/2"))
	}
	address := (*m)[o.ToAddress]
	for _, cell := range []string{"H4", "H16"} {
		f.SetCellValue(sheetName, cell, address.String())
	}

	// CreateAllocate(f, c)
	a := new(Allocation)
	if err := c.Bind(&o); err != nil {
		msg := err.Error()
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": msg, "error": err})
		return
	}
	fmt.Printf("Allocation: %#v\n", a)

	sheetName = "入力画面"
	// 要求番号
	reqNo, err := getRequestNo(a.Section)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "error": err})
		return
	}
	f.SetCellValue(sheetName, "F2", reqNo.Base)
	// 要求年月日
	f.SetCellValue(sheetName, "F3", time.Now().Format(LAYOUT))
	f.SetCellValue(sheetName, "F4", a.Section)
	// 輸送便の別
	s := a.Type
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
	f.SetCellValue(sheetName, "F7", a.Order)
	// 輸送区間
	from, _ := f.GetCellValue(sheetName, "F8")
	to := "適当な宛先"
	f.SetCellValue(sheetName, "F8", from+to)
	// 送り先
	addresses := new(api.AddressMap)
	if err := api.UnmarshalJSON(addresses, api.ADDRESSFILE); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"msg": err.Error(), "error": err})
		return
	}
	f.SetCellValue(sheetName, "F13", strings.Join((*addresses)[to], "\n"))
	// t数 台車 機能
	t := fmt.Sprintf("%st%s(%s)", strconv.Itoa(a.T), a.Cartype, a.Function)
	f.SetCellValue(sheetName, "F6", t)
	// 積込/到着作業月日/時刻
	x := map[string]interface{}{
		"F9":  a.Load.Date.Format(LAYOUT),
		"F10": fmt.Sprintf("%d時%d分", a.Load.Hour, a.Load.Minute),
		"F11": a.Arrive.Date.Format(LAYOUT),
		"F12": fmt.Sprintf("%d時%d分", a.Arrive.Hour, a.Arrive.Minute),
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
	if a.Insulance != "契約済み" {
		x = map[string]interface{}{
			"F18": `☐要　☑不要`,         // 保険
			"F19": a.InsulancePrice, // 保険額
		}
	}
	for cell, value := range x {
		if err = f.SetCellValue(sheetName, cell, value); err != nil {
			fmt.Println(err)
		}
	}
	// 寸法質量
	l := len(a.Package)
	if l < MAXLINE { // 3行までなら配車要求票に記載
		p := a.Size.Compile()
		fmt.Println(p)
		x := map[string]interface{}{
			"F15": a.Size.ToString(), // 重量・長さなど
			"F16": p.ToString(),      // 荷姿(個数)
			"F17": a.Size.Sum(),      // 総個数
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
			"F17": a.Size.Sum(),
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
				"F" + n: a.Package[i],
				"G" + n: fmt.Sprintf("%dx%dx%d", a.Width[i], a.Length[i], a.Hight[i]),
				"J" + n: a.Mass[i],
				"K" + n: a.Method[i],
				"L" + n: a.Quantity[i],
			}
			for cell, value := range x {
				if err = f.SetCellValue(sheetName, cell, value); err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	downloadFile("外装ラベル.xlsx", f, c)
}

func downloadFile(fs string, f *excelize.File, c *gin.Context) {
	// Save file
	filepath := "./result.xlsx"
	if err := f.SaveAs(filepath); err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(filepath)

	// Download file
	buf, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Stream(func(w io.Writer) bool {
		ar := zip.NewWriter(w)
		defer ar.Close()
		c.Writer.Header().Set("Content-Disposition", "attachmnt; filename=download.zip")
		f, _ := ar.Create(fs)
		io.Copy(f, buf)
		return false
	})
}
