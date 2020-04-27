package config

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultRegion        = "us-east-1"
	defaultDynamoDBTable = "streaming-users-online"
	defaultHTTPHost      = ":8888"
)

var (
	configETCPath            = "/etc/ws-message-dispatcher/"
	configLocalPath          = "./config"
	configType               = "yaml"
	configFileName           = "ws-message-dispatcher"
	configDynamoRegion       = "dynamodb.region"
	configDynamoTableName    = "dynamodb.table"
	configLambdaRegion       = "lambda.region"
	configLambdaFunctionName = "lambda.function"
	configServiceHost        = "http.host"

	errMissingConfiguration = errors.New("missing configuration")
)

type dynamoConfig struct {
	Region string
	Table  string
}

type lambdaConfig struct {
	Region   string
	Function string
}

type http struct {
	Host string
}

// Config holds service config
type Config struct {
	Dynamo  dynamoConfig
	Lambda  lambdaConfig
	Service http
}

// ReadConfig reads config service
func ReadConfig() (Config, error) {
	viper.AddConfigPath(configETCPath)
	viper.AddConfigPath(configLocalPath)

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configType)

	viper.SetDefault(configDynamoRegion, defaultRegion)
	viper.SetDefault(configDynamoTableName, defaultDynamoDBTable)
	viper.SetDefault(configServiceHost, defaultHTTPHost)

	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"config_file": configFileName,
		}).Error("unable to read config file")
		return Config{}, err
	}

	conf := Config{}

	conf.Dynamo.Region = viper.GetString(configDynamoRegion)
	conf.Dynamo.Table = viper.GetString(configDynamoTableName)
	conf.Lambda.Region = viper.GetString(configLambdaRegion)
	conf.Lambda.Function = viper.GetString(configLambdaFunctionName)
	conf.Service.Host = viper.GetString(configServiceHost)

	if len(conf.Lambda.Function) == 0 {
		log.Error("lambda function does not set")
		return conf, errMissingConfiguration
	}

	log.WithFields(log.Fields{
		"dynamo-Region":   conf.Dynamo.Region,
		"dynamo-table":    conf.Dynamo.Table,
		"lambda-Region":   conf.Lambda.Region,
		"lambda-function": conf.Lambda.Function,
		"http-host":       conf.Service.Host,
	}).Info("config read")

	return conf, nil
}
