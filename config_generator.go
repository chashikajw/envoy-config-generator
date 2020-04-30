package main

import (

	swgger "envoy-config-generator/swaggerOperator"
	"fmt"
)

//package envoy_config_generator

func main() {
	mg := swgger.GenerateMgwSwagger()
	fmt.Println(mg)
}

