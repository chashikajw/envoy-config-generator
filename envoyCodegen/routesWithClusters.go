package envoyCodegen

import (

	"envoy-config-generator/models/apiDefinition"
	"errors"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	v2route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/golang/protobuf/ptypes"
	"strings"
	"time"
)
func createRoutesWithClusters(mgwSwagger  apiDefinition.MgwSwagger) ([]*v2route.Route, []*v2.Cluster) {
	var (
		routes []*v2route.Route
		clusters []*v2.Cluster
		endpoint apiDefinition.Endpoint
		apiLevelEndpoint apiDefinition.Endpoint
		cluster v2.Cluster
		apilevelCluster v2.Cluster
		clusterName string
		address core.Address
		cluster_ref string

	)

	if(isProductionEndpointsAvailable(mgwSwagger.VendorExtensible)) {
		apiLevelEndpoint  = getProductionEndpoint(mgwSwagger.VendorExtensible)
		apilevelAddress := createAddress(apiLevelEndpoint.Url[0],config.API_PORT)
		apiLevelClusterName := "cluster_" + strings.Replace(mgwSwagger.Title, " ", "", -1) +  mgwSwagger.Version
		apilevelCluster = createCluster(apilevelAddress,apiLevelClusterName)
		clusters = append(clusters, &apilevelCluster)


	} else{
		errors.New("Producton endpoints are not defined")
	}

	for ind,resource:= range mgwSwagger.Resources {

		//resource level check
		if(isProductionEndpointsAvailable(resource.VendorExtensible)) {
			endpoint  = getProductionEndpoint(resource.VendorExtensible)
			address = createAddress(endpoint.Url[0],config.API_PORT)
			clusterName = "cluster_" + strings.Replace(resource.ID, " ", "", -1) + string(ind)
			cluster = createCluster(address,clusterName)
			clusters = append(clusters, &cluster)

			cluster_ref = cluster.GetName()

			//API level check
		} else if(isProductionEndpointsAvailable(mgwSwagger.VendorExtensible)){
			endpoint = apiLevelEndpoint
			cluster_ref = apilevelCluster.GetName()

		} else {
			errors.New("Producton endpoints are not defined")
		}

		route := createRoute(endpoint.Url[0], resource.Context,cluster_ref)
		routes = append(routes, &route)
	}

	return routes, clusters

}

func createCluster(address core.Address,clusterName string) v2.Cluster {

	h := &address
	cluster := v2.Cluster{
		Name:                 clusterName,
		ConnectTimeout:       ptypes.DurationProto(2 * time.Second),
		ClusterDiscoveryType: &v2.Cluster_Type{Type: v2.Cluster_LOGICAL_DNS},
		DnsLookupFamily:      v2.Cluster_V4_ONLY,
		LbPolicy:             v2.Cluster_ROUND_ROBIN,
		Hosts:                []*core.Address{h},
	}

	return cluster
}

func createRoute(HostUrl string, resourcePath string, clusterName string) v2route.Route {
	var targetRegex = "/"
	route := v2route.Route{
		Match:&v2route.RouteMatch{
			PathSpecifier: &v2route.RouteMatch_Regex{
				Regex: targetRegex,
			},
		},
		Action: &v2route.Route_Route{
			Route: &v2route.RouteAction{
				HostRewriteSpecifier: &v2route.RouteAction_HostRewrite{
					HostRewrite: HostUrl,
				},
				ClusterSpecifier: &v2route.RouteAction_Cluster{
					Cluster: clusterName,
				},
				PrefixRewrite: resourcePath,
			},
		},
		Metadata:     nil,
	}
	return route
}

