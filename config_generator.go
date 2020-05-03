package main
//package envoy_config_generator

import (
	e "envoy-config-generator/envoyCodegen"
	"envoy-config-generator/models/envoy"
	swgger "envoy-config-generator/swaggerOperator"
	"fmt"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"strings"
)



func main() {

	fmt.Println(GetSandboxsources())

}

func GetProductionsources() ([]types.Resource, []types.Resource, []types.Resource, []types.Resource) {
	mgwSwagger := swgger.GenerateMgwSwagger()
	//fmt.Println(mgwSwagger)

	routesP, clustersP, _, _ := e.CreateRoutesWithClusters(mgwSwagger)

	vHost_NameP := "serviceProd_" + strings.Replace(mgwSwagger.Title, " ", "", -1) +  mgwSwagger.Version

	vHostP, err := e.CreateVirtualHost(vHost_NameP,routesP)

	fmt.Println(err)

	listenerNameP := "listenerProd_1"
	routeConfigNameP := "routeProd_" + strings.Replace(mgwSwagger.Title, " ", "", -1) +  mgwSwagger.Version

	listnerProd := e.CreateListener(listenerNameP,routeConfigNameP, vHostP)

	envoyNodeProd := new(envoy.EnvoyNode)
	envoyNodeProd.SetListener(&listnerProd)
	envoyNodeProd.SetClusters(clustersP)
	envoyNodeProd.SetRoutes(routesP)
	fmt.Println(routesP)
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++")
	return envoyNodeProd.GetSources()
}

func GetSandboxsources() ([]types.Resource, []types.Resource, []types.Resource, []types.Resource) {
	mgwSwagger := swgger.GenerateMgwSwagger()
	//fmt.Println(mgwSwagger)

	_, _, routesS, clustersS := e.CreateRoutesWithClusters(mgwSwagger)
	if(routesS == nil) {
		return nil, nil, nil, nil
	}

	vHost_NameS := "serviceSand_" + strings.Replace(mgwSwagger.Title, " ", "", -1) +  mgwSwagger.Version

	vHostS, err := e.CreateVirtualHost(vHost_NameS,routesS)

	fmt.Println(err)

	listenerNameS := "listenerSand_1"
	routeConfigNameS := "routeSand_" + strings.Replace(mgwSwagger.Title, " ", "", -1) +  mgwSwagger.Version

	listnerSand := e.CreateListener(listenerNameS,routeConfigNameS, vHostS)

	envoyNodeSand := new(envoy.EnvoyNode)
	envoyNodeSand.SetListener(&listnerSand)
	envoyNodeSand.SetClusters(clustersS)
	envoyNodeSand.SetRoutes(routesS)
	fmt.Println(routesS)
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++")

	return envoyNodeSand.GetSources()
}