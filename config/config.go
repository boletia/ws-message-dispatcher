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

	envConfigDynamoRegion       = "DYNAMODB_REGION"
	envConfigDynamoTableName    = "DYNAMODB_TABLE"
	envConfigLambdaRegion       = "LAMBDA_REGION"
	envConfigLambdaFunctionName = "LAMBDA_FUNCTION"
	envConfigServiceHost        = "HTTP_HOST"

	errMissingConfiguration   = errors.New("missing configuration")
	errUnableToReadConfigFile = errors.New("unable to read config file")
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

// Read reads config service
func Read() (Config, error) {
	conf := Config{}
	var err error

	err = readFromFile(&conf)
	if err != nil && err != errUnableToReadConfigFile {
		return Config{}, err
	}

	if err == nil {
		log.WithFields(log.Fields{
			"dynamo-Region":   conf.Dynamo.Region,
			"dynamo-table":    conf.Dynamo.Table,
			"lambda-Region":   conf.Lambda.Region,
			"lambda-function": conf.Lambda.Function,
			"http-host":       conf.Service.Host,
		}).Info("config read from file")

		return conf, nil
	}

	log.Warn("unable to read config file, trying with env vars")
	err = readFromEnv(&conf)
	if err != nil {
		return Config{}, err
	}

	log.WithFields(log.Fields{
		"dynamo-Region":   conf.Dynamo.Region,
		"dynamo-table":    conf.Dynamo.Table,
		"lambda-Region":   conf.Lambda.Region,
		"lambda-function": conf.Lambda.Function,
		"http-host":       conf.Service.Host,
	}).Info("config read from envs")

	return conf, nil
}

func readFromEnv(conf *Config) error {
	configVars := map[string]string{
		envConfigDynamoRegion:       "",
		envConfigDynamoTableName:    "",
		envConfigLambdaRegion:       "",
		envConfigLambdaFunctionName: "",
		envConfigServiceHost:        "",
	}

	for key := range configVars {
		viper.BindEnv(key)

		if len(viper.GetString(key)) == 0 {
			log.WithFields(log.Fields{
				"variable": key,
				"value":    viper.GetString(key),
			}).Error("reading environment config")
			return errMissingConfiguration
		}

		configVars[key] = viper.GetString(key)
	}

	conf.Dynamo.Region = configVars[envConfigDynamoRegion]
	conf.Dynamo.Table = configVars[envConfigDynamoTableName]
	conf.Lambda.Region = configVars[envConfigLambdaRegion]
	conf.Lambda.Function = configVars[envConfigLambdaFunctionName]
	conf.Service.Host = configVars[envConfigServiceHost]

	return nil
}

func readFromFile(conf *Config) error {
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
		return errUnableToReadConfigFile
	}

	conf.Dynamo.Region = viper.GetString(configDynamoRegion)
	conf.Dynamo.Table = viper.GetString(configDynamoTableName)
	conf.Lambda.Region = viper.GetString(configLambdaRegion)
	conf.Lambda.Function = viper.GetString(configLambdaFunctionName)
	conf.Service.Host = viper.GetString(configServiceHost)

	if len(conf.Lambda.Function) == 0 {
		log.Error("lambda function does not set")
		return errMissingConfiguration
	}

	return nil
}
