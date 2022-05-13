package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/u1and0/schd/cmd/api"
	"github.com/u1and0/schd/cmd/paper"
)

const (
	// VERSION : schd version
	VERSION = "v0.1.0"
)

var (
	showVersion bool
)

func main() {
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.Parse()
	if showVersion {
		fmt.Println("schd version", VERSION)
		os.Exit(0) // Exit with version info
	}
	// Router
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("template/*.tmpl")

	// API
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	v1 := r.Group("api/v1")
	{
		d := v1.Group("/data")
		{
			d.GET("/", api.All)
			d.GET("/:id", api.Get)
			d.POST("/add", api.Post)
			d.DELETE("/:id", api.Delete)
			d.PUT("/:id", api.Put)
			d.GET("/list", api.List)
			d.GET("/address", api.FetchAddress)
			d.GET("/allocate", api.FetchAllocate)
			d.GET("/allocate/:id", api.FetchAllocateID)
			d.GET("/allocate/search", api.SearchAllocate)
		}

		v := v1.Group("/view")
		{
			v.GET("/list", api.View) // 日付リスト
			v.GET("/:id", api.Show)
			v.GET("/add/form", api.CreateForm)
			v.POST("/add", api.Create)
			v.GET("/:id/update/form", api.RefreshForm)
			v.POST("/:id/update", api.Refresh)
			v.GET("/:id/delete", api.Remove)
		}

		p := v1.Group("paper")
		{
			l := p.Group("label")
			{
				l.GET("/form", paper.CreateForm)
				l.POST("/post", paper.Create)
			}
			a := p.Group("allocate")
			{
				a.GET("/form", paper.CreateAllocateForm)
				a.POST("/post", paper.CreateAllocate)
			}
		}
	}

	r.Run()
}
