package envoyCodegen

import (
	"envoy-config-generator/config"
	"envoy-config-generator/models/apiDefinition"
	c "envoy-config-generator/constants"
	s "envoy-config-generator/swaggerOperator"
	"errors"
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	v2route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/golang/protobuf/ptypes"
	"strings"
	"time"
)
func CreateRoutesWithClusters(mgwSwagger  apiDefinition.MgwSwagger) ([]*v2route.Route, []*v2.Cluster, []*v2route.Route, []*v2.Cluster) {
	var (
		routesP []*v2route.Route
		clustersP []*v2.Cluster
		endpointP apiDefinition.Endpoint
		apiLevelEndpointP apiDefinition.Endpoint
		clusterP v2.Cluster
		apilevelClusterP v2.Cluster
		clusterNameP string
		addressP core.Address
		cluster_refP string

		routesS []*v2route.Route
		clustersS []*v2.Cluster
		endpointS apiDefinition.Endpoint
		apiLevelEndpointS apiDefinition.Endpoint
		clusterS v2.Cluster
		apilevelClusterS v2.Cluster
		clusterNameS string
		addressS core.Address
		cluster_refS string

	)

	if(s.IsSandboxEndpointsAvailable(mgwSwagger.VendorExtensible)) {
		apiLevelEndpointS  = s.GetEndpoints(mgwSwagger.VendorExtensible, c.SANDBOX_ENDPOINTS)
		apilevelAddressS := createAddress(apiLevelEndpointS.Url[0],config.API_PORT)
		apiLevelClusterNameS := "clusterSand_" + strings.Replace(mgwSwagger.Title, " ", "", -1) +  mgwSwagger.Version
		apilevelClusterS = createCluster(apilevelAddressS,apiLevelClusterNameS)
		clustersS = append(clustersS, &apilevelClusterS)
	}

	if(s.IsProductionEndpointsAvailable(mgwSwagger.VendorExtensible)) {
		apiLevelEndpointP  = s.GetEndpoints(mgwSwagger.VendorExtensible, c.PRODUCTION_ENDPOINTS)
		apilevelAddressP := createAddress(apiLevelEndpointP.Url[0],config.API_PORT)
		apiLevelClusterNameP := "clusterProd_" + strings.Replace(mgwSwagger.Title, " ", "", -1) +  mgwSwagger.Version
		apilevelClusterP = createCluster(apilevelAddressP,apiLevelClusterNameP)
		clustersP = append(clustersP, &apilevelClusterP)


	} else{
		errors.New("Producton endpoints are not defined")
	}

	for ind,resource:= range mgwSwagger.Resources {

		//resource level check sandbox endpoints
		if(s.IsSandboxEndpointsAvailable(resource.VendorExtensible)) {
			endpointS  = s.GetEndpoints(resource.VendorExtensible, c.SANDBOX_ENDPOINTS)
			addressS = createAddress(endpointS.Url[0],config.API_PORT)
			clusterNameS = "clusterSand_" + strings.Replace(resource.ID, " ", "", -1) + string(ind)
			clusterS = createCluster(addressS,clusterNameS)
			clustersS = append(clustersS, &clusterS)

			cluster_refS = clusterS.GetName()

			//sandbox endpoints
			routeS := createRoute(endpointS.Url[0], resource.Context,cluster_refS)
			routesS = append(routesS, &routeS)

			//API level check
		} else if(s.IsSandboxEndpointsAvailable(mgwSwagger.VendorExtensible)) {
			endpointS = apiLevelEndpointS
			cluster_refS = apilevelClusterS.GetName()

			//sandbox endpoints
			routeS := createRoute(endpointS.Url[0], resource.Context,cluster_refS)
			routesS = append(routesS, &routeS)
		}

		//resource level check
		if(s.IsProductionEndpointsAvailable(resource.VendorExtensible)) {
			endpointP  = s.GetEndpoints(resource.VendorExtensible, c.PRODUCTION_ENDPOINTS)
			addressP = createAddress(endpointP.Url[0],config.API_PORT)
			clusterNameP = "clusterProd_" + strings.Replace(resource.ID, " ", "", -1) + string(ind)
			clusterP = createCluster(addressP,clusterNameP)
			clustersP = append(clustersP, &clusterP)

			cluster_refP = clusterP.GetName()

			//production endpoints
			routeP := createRoute(endpointP.Url[0], resource.Context,cluster_refP)
			routesP = append(routesP, &routeP)

			//API level check
		} else if(s.IsProductionEndpointsAvailable(mgwSwagger.VendorExtensible)){
			endpointP = apiLevelEndpointP
			cluster_refP = apilevelClusterP.GetName()

			//production endpoints
			routeP := createRoute(endpointP.Url[0], resource.Context,cluster_refP)
			routesP = append(routesP, &routeP)

		} else {
			errors.New("Producton endpoints are not defined")
		}


	}

	return routesP, clustersP,  routesS, clustersS

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

