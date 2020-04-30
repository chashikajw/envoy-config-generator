package main

import (

	"envoy-config-generator/utills"
	"envoy-config-generator/models/apiDefinition"

	"encoding/json"

	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/spec"
	"io/ioutil"
	"os"



)

//package envoy_config_generator

func main() {
	mg := GenerateMgwSwagger()
	fmt.Println(mg)
}

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