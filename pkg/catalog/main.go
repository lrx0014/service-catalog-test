package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	client "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/svcat"
	catalog "github.com/kubernetes-incubator/service-catalog/pkg/svcat/service-catalog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
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
		authorized.GET("/v1/catalog/test", addBroker)

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

func addBroker(c *gin.Context) {
	cfg, err1 := rest.InClusterConfig()
	if err1 != nil {
		result := make(map[string]interface{})
		handleError(c, "GET", result, err1, "addBroker")
	}
	catalogClient, err2 := client.NewForConfig(cfg)
	if err2 != nil {
		result := make(map[string]interface{})
		handleError(c, "GET", result, err2, "addBroker")
	}
	k8sClient, err3 := kubernetes.NewForConfig(cfg)
	if err3 != nil {
		result := make(map[string]interface{})
		handleError(c, "GET", result, err3, "addBroker")
	}
	app, err4 := svcat.NewApp(k8sClient, catalogClient, "default")
	if err4 != nil {
		result := make(map[string]interface{})
		handleError(c, "GET", result, err4, "addBroker")
	}
	opt := &catalog.RegisterOptions{
		BasicSecret: "test",
		SkipTLS:     true,
		Namespace:   "default",
	}
	scope := &catalog.ScopeOptions{
		Scope: "cluster",
	}
	broker, err5 := app.Register("fake-broker", "http://fake-broker.io", opt, scope)
	if err5 != nil {
		result := make(map[string]interface{})
		handleError(c, "GET", result, err5, "addBroker")
	}
	status := broker.GetStatus()
	log.Println(status)
	result := make(map[string]interface{})
	handleSuccess(c, result, "addBroker")
}

func handleSuccess(c *gin.Context, result interface{}, httpmethod string) {
	c.JSON(200, result)
}

func handleError(c *gin.Context, request string, result map[string]interface{}, err error, httpmethod string) {
	response := make(map[string]interface{})
	if err.Error() == "exist" {
		c.JSON(409, response)
	} else if err.Error() == "noexist" {
		c.JSON(410, response)
	} else {
		response["error"] = request
		response["description"] = fmt.Sprintf("Error: %v", err)
		c.JSON(500, response)
	}
}
