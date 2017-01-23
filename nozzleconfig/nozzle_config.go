package nozzleconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type NozzleConfig struct {
	UAAURL                 string
	Client                 string
	ClientSecret           string
	TrafficControllerURL   string
	FirehoseSubscriptionID string
	InfluxDbUrl            string
	InfluxDbDatabase       string
	InfluxDbUser           string
	InfluxDbPassword       string
	InfluxDbSslSkipVerify  bool
	FlushDurationSeconds   uint32
	SsLSkipVerify          bool
	MetricPrefix           string
	Deployment             string
	DeploymentFilter       string
	EventFilter            string
	DisableAccessControl   bool
	IdleTimeoutSeconds     uint32
	AppInfoApiUrl	       string
}

func Parse(configPath string) (*NozzleConfig, error) {
	configBytes, err := ioutil.ReadFile(configPath)
	var config NozzleConfig
	if err != nil {
		return nil, fmt.Errorf("Can not read config file [%s]: %s", configPath, err)
	}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("Can not parse config file %s: %s", configPath, err)
	}

	overrideWithEnvVar("NOZZLE_UAAURL", &config.UAAURL)
	overrideWithEnvVar("NOZZLE_CLIENT", &config.Client)
	overrideWithEnvVar("NOZZLE_CLIENT_SECRET", &config.ClientSecret)
	overrideWithEnvVar("NOZZLE_TRAFFICCONTROLLERURL", &config.TrafficControllerURL)
	overrideWithEnvVar("NOZZLE_FIREHOSESUBSCRIPTIONID", &config.FirehoseSubscriptionID)
	overrideWithEnvVar("NOZZLE_INFLUXDB_URL", &config.InfluxDbUrl)
	overrideWithEnvVar("NOZZLE_INFLUXDB_DATABASE", &config.InfluxDbDatabase)
	overrideWithEnvVar("NOZZLE_INFLUXDB_USER", &config.InfluxDbUser)
	overrideWithEnvVar("NOZZLE_INFLUXDB_PASSWORD", &config.InfluxDbPassword)
	overrideWithEnvBool("NOZZLE_INFLUXDB_SSL_SKIPVERIFY", &config.InfluxDbSslSkipVerify)
	overrideWithEnvVar("NOZZLE_METRICPREFIX", &config.MetricPrefix)
	overrideWithEnvVar("NOZZLE_DEPLOYMENT", &config.Deployment)
	overrideWithEnvVar("NOZZLE_DEPLOYMENT_FILTER", &config.DeploymentFilter)
	overrideWithEnvVar("NOZZLE_EVENT_FILTER", &config.EventFilter)

	overrideWithEnvUint32("NOZZLE_FLUSHDURATIONSECONDS", &config.FlushDurationSeconds)

	overrideWithEnvBool("NOZZLE_SSL_SKIPVERIFY", &config.SsLSkipVerify)
	overrideWithEnvBool("NOZZLE_DISABLEACCESSCONTROL", &config.DisableAccessControl)
	overrideWithEnvUint32("NOZZLE_IDLETIMEOUTSECONDS", &config.IdleTimeoutSeconds)
	overrideWithEnvVar("NOZZLE_APP_API_URL", &config.AppInfoApiUrl)

	return &config, nil
}

func overrideWithEnvVar(name string, value *string) {
	envValue := os.Getenv(name)
	if envValue != "" {
		*value = envValue
	}
}

func overrideWithEnvUint32(name string, value *uint32) {
	envValue := os.Getenv(name)
	if envValue != "" {
		tmpValue, err := strconv.Atoi(envValue)
		if err != nil {
			panic(err)
		}
		*value = uint32(tmpValue)
	}
}

func overrideWithEnvBool(name string, value *bool) {
	envValue := os.Getenv(name)
	if envValue != "" {
		var err error
		*value, err = strconv.ParseBool(envValue)
		if err != nil {
			panic(err)
		}
	}
}
