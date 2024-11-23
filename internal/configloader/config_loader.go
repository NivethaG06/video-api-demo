package configloader

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

func LoadConfig(filePath string) (*Configurations, error) {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %v", err)
	}

	var config Configurations
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML file: %v", err)
	}

	return &config, nil
}

func SetEnvVars(config interface{}, prefix string) {
	v := reflect.ValueOf(config).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fieldName := field.Name
		fieldValue := v.Field(i).Interface()

		if field.Type.Kind() == reflect.Struct {
			SetEnvVars(v.Field(i).Addr().Interface(), prefix+"_"+strings.ToUpper(fieldName))
		} else {
			envVarName := strings.ToUpper(prefix + "_" + fieldName)
			err := os.Setenv(envVarName, fmt.Sprintf("%v", fieldValue))
			if err != nil {
				log.Printf("Failed to set environment variable %s: %v", envVarName, err)
			}
		}
	}
}
