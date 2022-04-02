package main

import (
	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/api"
)

func main() {
	// Router
	r := gin.Default()
	r.LoadHTMLGlob("template/*.tmpl")

	// API
	v1 := r.Group("api/v1")
	{
		d := v1.Group("/data")
		{
			d.GET("/all", api.Index)
			d.GET("/:id", api.Get)
			d.POST("/add", api.Post)
			d.DELETE("/:id", api.Delete)
			d.PUT("/:id", api.Put)
			d.GET("/list", api.List)
		}

		v := v1.Group("view")
		{
			v.GET("/list", api.View) // 日付リスト
			v.GET("/:id", api.Show)
			v.GET("/add/form", api.CreateForm)
			v.POST("/add", api.Create)
			v.PUT("/update/:id", api.Update)
			// v.GET("/update/form/:id", api.UpdateForm)
			v.GET("/delete/:id", api.Remove)
		}
	}

	r.Run()
}
