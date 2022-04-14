package ctrl

import (
	"strconv"
	"time"
)

/*
form入力
	http://...api/list/2204
	日付iter
	2022-04-01 ~ 2022-04-30 までの日付リストを取得
*/

func ToMonth(s string) (int, int, error) {
	ys := "20" + s[:2]
	ms := s[2:]
	y, err := strconv.Atoi(ys)
	m, err := strconv.Atoi(ms)
	return y, m, err
}
func Iter(st, en time.Time) {

}
