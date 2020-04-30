package swaggerOperator

import (
	"envoy-config-generator/constants"
	"envoy-config-generator/utills"
	"envoy-config-generator/models/apiDefinition"

	"encoding/json"

	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/spec"
	"io/ioutil"
	"os"
)

func GenerateMgwSwagger() apiDefinition.MgwSwagger {
	var mgwSwagger  apiDefinition.MgwSwagger

	// Open our jsonFile
	openApif, err := os.Open("test.yaml")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened open api file")
	// defer the closing of our jsonFile so that we can parse it later on
	defer openApif.Close()

	// read our opened jsonFile as a byte array.
	jsn, _ := ioutil.ReadAll(openApif)

	apiJsn, err := utills.ToJSON(jsn)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	swaggerVerison, err := utills.FindSwaggerVersion(apiJsn)

	if(swaggerVerison == "2") {
		//map json to struct
		var ApiData spec.Swagger
		err = json.Unmarshal(apiJsn, &ApiData)
		if err != nil {
			fmt.Printf("openAPI unmarsheliing err: %v\n", err)
		} else {
			mgwSwagger.SetInfoSwagger(ApiData)
		}

	}else if(swaggerVerison == "3") {
		//map json to struct
		var ApiData openapi3.Swagger
		err = json.Unmarshal(apiJsn, &ApiData)
		if err != nil {
			fmt.Printf("openAPI unmarsheliing err: %v\n", err)
		} else {
			mgwSwagger.SetInfoOpenApi(ApiData)
		}
	}
	return mgwSwagger
}

func IsProductionEndpointsAvailable(vendorExtensible map[string]interface{}) bool{
	if _, found := vendorExtensible[constants.PRODUCTION_ENDPOINTS]; found {
		return true
	} else {
		return false
	}
}



func GetProductionEndpoint(vendorExtensible map[string]interface{}) apiDefinition.Endpoint{
	var productionEndpoint apiDefinition.Endpoint
	if y, found := vendorExtensible["x-wso2-production-endpoints"]; found {
		if val, ok := y.(map[string]interface{}); ok {

			for ind,val := range val {
				//fmt.Println(ind, val)
				if(ind == "type") {
					productionEndpoint.UrlType = val.(string)
				} else if (ind == "urls") {
					ainterface := val.([]interface {})
					urls := make([]string, len(ainterface))
					for i, v := range ainterface {
						urls[i] = v.(string)
					}
					productionEndpoint.Url = urls
				}

			}

		} else {
			fmt.Println("production endpoint is not having a correct structure")
		}

	} else {
		fmt.Println("couldn't find a basepath extenstion")
	}


	return productionEndpoint
}
