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
		"today":   time.Now().Format(LAYOUT),
		"section": (*section),
		"hours":   hours,
		"minutes": minutes,
	})
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
