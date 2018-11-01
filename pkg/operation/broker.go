package operation

import (
	"log"

	v1beta1 "github.com/kubernetes-incubator/service-catalog/pkg/apis/servicecatalog/v1beta1"
	svcat "github.com/kubernetes-incubator/service-catalog/pkg/svcat"
	catalog "github.com/kubernetes-incubator/service-catalog/pkg/svcat/service-catalog"
)

// BrokerInterface define the Broker Methods
type BrokerInterface interface {
	AddBroker(client *svcat.App) (status v1beta1.CommonServiceBrokerStatus)
	GetBroker(client *svcat.App, name string) (info *v1beta1.ClusterServiceBroker)
}

// AddBroker add a broker
func AddBroker(client *svcat.App) (status v1beta1.CommonServiceBrokerStatus) {
	opt := &catalog.RegisterOptions{
		BasicSecret: "test",
		SkipTLS:     true,
		Namespace:   "default",
	}
	scope := &catalog.ScopeOptions{
		Scope: "cluster",
	}
	broker, err := client.Register("fake-broker", "http://fake-broker.io", opt, scope)
	if err != nil {
		log.Println(err)
	}
	status = broker.GetStatus()
	log.Println(status)
	return status
}

// GetBroker describe a broker's info
func GetBroker(client *svcat.App, name string) (info *v1beta1.ClusterServiceBroker) {
	broker, err := client.RetrieveBroker(name)
	if err != nil {
		log.Println(err)
	}
	return broker
}
