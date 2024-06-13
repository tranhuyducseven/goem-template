package openobserve

import (
	"fmt"
)

type OpenObserveConfig struct {
	Endpoint    string
	Credential  string
	ServiceName string
	Env         string
	OrgName     string

	logEndpoint string
}

var config OpenObserveConfig

func Init(conf OpenObserveConfig) {
	config = conf
	if config.OrgName == "" {
		config.OrgName = "default"
	}
	config.logEndpoint = fmt.Sprintf("%s/api/%s/%s/_json", config.Endpoint, config.OrgName, config.ServiceName)
}

func GetConfig() OpenObserveConfig {
	return config
}
