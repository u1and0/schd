package api

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ctrl"
)

type (
	// Searchers : slice of Searcher
	Searchers []Searcher
	// Searcher : struct for Searcher
	Searcher struct {
		ID   ctrl.ID   `json:"id"`
		Body string    `json:"body"`
		Date time.Time `json:"date"`
	}
)

// FetchAllocateList : returns Searcher list
func FetchAllocateList(c *gin.Context) {
	searchers := allocations.Concat()
	sort.Sort(searchers)
	c.IndentedJSON(http.StatusOK, searchers)
}

// SearchAllocate : allocate info search from query
func SearchAllocate(c *gin.Context) {
	result := make(Allocations, len(allocations))
	q := c.Query("q")
	searchers := allocations.Concat()
	keywd := strings.Split(q, " ")
	for _, s := range searchers {
		// sにkeywdが全て含まれていたらresultに加える
		// 順序関係なし
		if ContainsAll(s.Body, keywd...) {
			result[s.ID] = allocations[s.ID]
		}
	}
	c.IndentedJSON(http.StatusOK, result)
}

// ContainsAll : すべて含まれていたらtrue
func ContainsAll(s string, keywords ...string) bool {
	for _, k := range keywords {
		if !strings.Contains(s, k) {
			return false
		}
	}
	return true
}

func (s Searchers) Len() int           { return len(s) }
func (s Searchers) Less(i, j int) bool { return s[i].Date.After(s[j].Date) }
func (s Searchers) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
