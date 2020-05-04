package main
//package envoy_config_generator

import (
	e "envoy-config-generator/envoyCodegen"
	"envoy-config-generator/models/envoy"
	swgger "envoy-config-generator/swaggerOperator"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"strings"
)



func main() {

	GetProductionsources()

}

func GetProductionsources() ([]types.Resource, []types.Resource, []types.Resource, []types.Resource) {
	mgwSwagger := swgger.GenerateMgwSwagger()
	//fmt.Println(mgwSwagger)

	routesP, clustersP,endpointsP, _, _,_ := e.CreateRoutesWithClusters(mgwSwagger)

	vHost_NameP := "serviceProd_" + strings.Replace(mgwSwagger.Title, " ", "", -1) +  mgwSwagger.Version

	vHostP, _ := e.CreateVirtualHost(vHost_NameP,routesP)

	//fmt.Println(err)

	listenerNameP := "listenerProd_1"
	routeConfigNameP := "routeProd_" + strings.Replace(mgwSwagger.Title, " ", "", -1) +  mgwSwagger.Version

	listnerProd := e.CreateListener(listenerNameP,routeConfigNameP, vHostP)

	envoyNodeProd := new(envoy.EnvoyNode)
	envoyNodeProd.SetListener(&listnerProd)
	envoyNodeProd.SetClusters(clustersP)
	envoyNodeProd.SetRoutes(routesP)
	envoyNodeProd.SetEndpoints(endpointsP)
	//fmt.Println(endpointsP)

	return envoyNodeProd.GetSources()
}

func GetSandboxsources() ([]types.Resource, []types.Resource, []types.Resource, []types.Resource) {
	mgwSwagger := swgger.GenerateMgwSwagger()
	//fmt.Println(mgwSwagger)

	_, _,_, routesS, clustersS, endpointsS := e.CreateRoutesWithClusters(mgwSwagger)
	if(routesS == nil) {
		return nil, nil, nil, nil
	}

	vHost_NameS := "serviceSand_" + strings.Replace(mgwSwagger.Title, " ", "", -1) +  mgwSwagger.Version

	vHostS, _ := e.CreateVirtualHost(vHost_NameS,routesS)

	//fmt.Println(err)

	listenerNameS := "listenerSand_1"
	routeConfigNameS := "routeSand_" + strings.Replace(mgwSwagger.Title, " ", "", -1) +  mgwSwagger.Version

	listnerSand := e.CreateListener(listenerNameS,routeConfigNameS, vHostS)

	envoyNodeSand := new(envoy.EnvoyNode)
	envoyNodeSand.SetListener(&listnerSand)
	envoyNodeSand.SetClusters(clustersS)
	envoyNodeSand.SetRoutes(routesS)
	envoyNodeSand.SetEndpoints(endpointsS)
	//fmt.Println(endpointsS)
	return envoyNodeSand.GetSources()
}