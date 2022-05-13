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
	for id, s := range allocator {
		if strings.Contains(s, q) {
			result[id] = allocations[id]
		}
	}
	c.IndentedJSON(http.StatusOK, result)
}
