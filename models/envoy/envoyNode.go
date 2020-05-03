package envoy

import (
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	v2route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

type EnvoyNode struct {
	listeners []types.Resource
	clusters  []types.Resource
	routes []types.Resource
	endpoints []types.Resource
}


func (envoy *EnvoyNode) SetListener(listener *v2.Listener) {
	envoy.listeners = []types.Resource{listener}

}

func (envoy *EnvoyNode) SetClusters(clusters []*v2.Cluster) {
	for _,clusterP := range clusters {
		envoy.clusters = append(envoy.clusters,clusterP )
	}
}

func (envoy *EnvoyNode) SetRoutes(routes []*v2route.Route) {
	for _,routesP := range routes {
		envoy.routes = append(envoy.routes,routesP)
	}
}

func (envoy *EnvoyNode) GetSources() ([]types.Resource, []types.Resource, []types.Resource, []types.Resource) {
	return envoy.listeners, envoy.clusters, envoy.routes, envoy.endpoints
}




