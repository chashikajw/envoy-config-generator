package envoyCodegen

import (
	config "envoy-config-generator/config"
	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	v2route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"

)

func CreateListener(listenerName string,routeConfigName string ,vHostP v2route.VirtualHost) (v2.Listener) {
	listenerAddress := &core.Address_SocketAddress{
		SocketAddress: &core.SocketAddress{
			Protocol: core.SocketAddress_TCP,
			Address:  config.LISTENER_ADDRESS,
			PortSpecifier: &core.SocketAddress_PortValue{
				PortValue: config.LISTENER_PORT,
			},
		},
	}


	listenerFilters := CreateListenerFilters(routeConfigName,vHostP )

	listener := v2.Listener{
		Name: listenerName,
		Address: &core.Address{
			Address: listenerAddress,
		},
		FilterChains: []*listener.FilterChain{{
			Filters: listenerFilters},
		},
	}

	return listener

}

func CreateListenerFilters(routeConfigName  string,vHost v2route.VirtualHost) ([]*listener.Filter) {
	var filters []*listener.Filter


	//set connection manager filter for production
	managerP := CreateConectionManagerFilter(vHost,routeConfigName)


	pbst, err := ptypes.MarshalAny(managerP)
	if err != nil {
		panic(err)
	}
	connectionManagerFilterP := listener.Filter{
		Name: wellknown.HTTPConnectionManager,
		ConfigType: &listener.Filter_TypedConfig{
			TypedConfig: pbst,
		},
	}



	//add filters
	filters = append(filters,&connectionManagerFilterP)
	return filters

}

func CreateConectionManagerFilter(vHost v2route.VirtualHost, routeConfigName string) *hcm.HttpConnectionManager {

	httpFilters := GetHttpFilters()

	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.HttpConnectionManager_AUTO,
		StatPrefix: config.MANAGER_STATPREFIX,
		RouteSpecifier: &hcm.HttpConnectionManager_RouteConfig{
			RouteConfig: &v2.RouteConfiguration{
				Name:         routeConfigName,
				VirtualHosts: []*v2route.VirtualHost{&vHost},
			},
		},
		HttpFilters: httpFilters,
	}
	return manager
}

func CreateVirtualHost(vHost_Name string,routes []*v2route.Route) (v2route.VirtualHost, error) {

	vHost_Domains :=  []string{"*"}

	virtual_host := v2route.VirtualHost{
		Name: vHost_Name,
		Domains: vHost_Domains,
		Routes: routes,
	}

	return virtual_host, nil
}

func createAddress(remoteHost string,Port uint32 ) core.Address {
	address := core.Address{Address: &core.Address_SocketAddress{
		SocketAddress: &core.SocketAddress{
			Address:  remoteHost,
			Protocol: core.SocketAddress_TCP,
			PortSpecifier: &core.SocketAddress_PortValue{
				PortValue: uint32(Port),
			},
		},
	}}
	return address
}

