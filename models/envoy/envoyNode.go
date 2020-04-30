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
	envoy.listeners = append(envoy.listeners, listener)
}

func (envoy *EnvoyNode) SetClusters(clusters []*v2.Cluster) {
	for _,cluster := range clusters {
		envoy.clusters = append(envoy.clusters,cluster )
	}
}

func (envoy *EnvoyNode) SetRoutes(routes []*v2route.Route) {
	for _,routes := range routes {
		envoy.routes = append(envoy.routes,routes)
	}
}

func (envoy *EnvoyNode) GetSources() ([]types.Resource, []types.Resource, []types.Resource, []types.Resource) {
	return envoy.listeners, envoy.clusters, envoy.routes, envoy.endpoints
}


