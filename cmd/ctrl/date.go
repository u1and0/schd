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

// ToMonth : string yymm to year yy and month mm
func ToMonth(s string) (int, int, error) {
	ys := "20" + s[:2]
	ms := s[2:]
	y, err := strconv.Atoi(ys)
	m, err := strconv.Atoi(ms)
	return y, m, err
}

// DayofFirstEnd : DayofFirstEnd(22,4) => 2022-04-01, 2022-04-30
func DayofFirstEnd(y, m int) (time.Time, time.Time) {
	first := time.Date(y, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
	end := first.AddDate(0, 1, -1)
	return first, end
}
