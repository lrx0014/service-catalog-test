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

var catalogApp *svcat.App

func main() {
	catalogApp = buildClient()
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
		authorized.POST("/v1/catalog/test", addBroker)
		authorized.GET("/v1/catalog/:broker", getBroker)

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

func buildClient() *svcat.App {
	cfg, err1 := rest.InClusterConfig()
	if err1 != nil {
		log.Println(err1)
		os.Exit(-1)
	}
	catalogClient, err2 := client.NewForConfig(cfg)
	if err2 != nil {
		log.Println(err2)
		os.Exit(-1)
	}
	k8sClient, err3 := kubernetes.NewForConfig(cfg)
	if err3 != nil {
		log.Println(err3)
		os.Exit(-1)
	}
	app, err4 := svcat.NewApp(k8sClient, catalogClient, "default")
	if err4 != nil {
		log.Println(err4)
		os.Exit(-1)
	}
	return app
}

func addBroker(c *gin.Context) {
	opt := &catalog.RegisterOptions{
		BasicSecret: "test",
		SkipTLS:     true,
		Namespace:   "default",
	}
	scope := &catalog.ScopeOptions{
		Scope: "cluster",
	}
	broker, err5 := catalogApp.Register("fake-broker", "http://fake-broker.io", opt, scope)
	if err5 != nil {
		result := make(map[string]interface{})
		handleError(c, "GET", result, err5, "addBroker")
	}
	status := broker.GetStatus()
	log.Println(status)
	result := make(map[string]interface{})
	handleSuccess(c, result, "addBroker")
}

func getBroker(c *gin.Context) {
	broker, err := catalogApp.RetrieveBroker(c.Param("broker"))
	if err != nil {
		result := make(map[string]interface{})
		handleError(c, "GET", result, err, "getBroker")
	}
	sta := broker.GetStatus()
	handleSuccess(c, sta, "getBroker")
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
