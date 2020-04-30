package main

import (

	swgger "envoy-config-generator/swaggerOperator"
	"fmt"
	e "envoy-config-generator/envoyCodegen"
	"envoy-config-generator/models/envoy"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

//package envoy_config_generator

func main() {

	fmt.Println(Getsources())

}

func Getsources() ([]types.Resource, []types.Resource, []types.Resource, []types.Resource) {
	mgwSwagger := swgger.GenerateMgwSwagger()
	fmt.Println(mgwSwagger)

	routes, clusters := e.CreateRoutesWithClusters(mgwSwagger)
	vHost, err := e.CreateVirtualHost(mgwSwagger,routes)
	fmt.Println(err)

	listenerName := "listener_1"
	listner := e.CreateListener(listenerName, mgwSwagger, vHost)
	envoyNode := new(envoy.EnvoyNode)
	envoyNode.SetListener(&listner)
	envoyNode.SetClusters(clusters)
	envoyNode.SetRoutes(routes)
	fmt.Println(envoyNode.GetSources())

	return envoyNode.GetSources()
}

