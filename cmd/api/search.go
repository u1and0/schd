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
	var contain bool
	result := make(Allocations, len(allocations))
	q := c.Query("q")
	allocator := allocations.Concat()
	keywd := strings.Split(q, " ")
	for id, s := range allocator {
		// keywdのうちすべて含んでいればtrue
		for _, k := range keywd {
			if !strings.Contains(s, k) {
				contain = false //bool 逆？
			}
			if contain {
				result[id] = allocations[id]
			}

		}
	}
	c.IndentedJSON(http.StatusOK, result)
}
