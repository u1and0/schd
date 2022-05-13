package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/ctrl"
)

type (
	// Allocator : struct for search
	Allocator map[ctrl.ID]string
)

// SearchAllocate : allocate info search from query
func SearchAllocate(c *gin.Context) {
	result := make(Allocations, len(allocations))
	q := c.Query("q")
	allocator := allocations.Concat()
	keywd := strings.Split(q, " ")
	for id, s := range allocator {
		// sにkeywdが全て含まれていたらresultに加える
		// 順序関係なし
		if ContainsAll(s, keywd...) {
			result[id] = allocations[id]
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
