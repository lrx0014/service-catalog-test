package main

import (
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/svcat"

	handler "github.com/lrx0014/service-catalog-test/pkg/handler"
	op "github.com/lrx0014/service-catalog-test/pkg/operation"
)

var catalogApp *svcat.App

func main() {
	catalogApp = op.NewClient()
	app := cli.NewApp()
	app.Name = "catalog-test"
	app.Usage = "Start the catalog components"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "",
			Usage: "config file used to run programmer",
		},
	}
	app.Action = func(c *cli.Context) {
		router := gin.Default()
		router.GET("/ping", func(c *gin.Context) {
			c.String(200, "pong")
		})
		authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
			"admin": "admin",
		}))
		//stack
		authorized.POST("/v1/catalog/test", handler.AddBroker)
		authorized.GET("/v1/catalog/:broker", handler.GetBroker)

		server := &http.Server{
			Addr:           ":9000",
			Handler:        router,
			ReadTimeout:    300 * time.Second,
			WriteTimeout:   300 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		server.ListenAndServe()
	}
	app.Run(os.Args)
}
