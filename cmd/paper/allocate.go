package paper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/api"
)

type Section map[string]string

// CreateForm : xlsxに転記するフォームの表示string
func CreateAllocateForm(c *gin.Context) {
	var section Section
	if err := api.UnmarshalJSON(section, api.SECTIONFILE); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "allocate_create.tmpl", gin.H{"section": section})
}
func CreateAllocate(c *gin.Context) {
}
