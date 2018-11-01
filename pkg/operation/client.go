package operation

import (
	"log"
	"os"

	client "github.com/kubernetes-incubator/service-catalog/pkg/client/clientset_generated/clientset"
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/svcat"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewClient create a catalog client
func NewClient() *svcat.App {
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
