package config

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultRegion                = "us-east-1"
	defaultDynamoUsersDBTable    = "streaming-users-online"
	defaultDynamoServersDBTable  = "chat-servers"
	defaultDynamoChatConfigTable = "streaming-dispatcher-config"
	defaultHTTPHost              = ":8888"
)

var (
	configETCPath                   = "/etc/ws-message-dispatcher/"
	configLocalPath                 = "./config"
	configType                      = "yaml"
	configFileName                  = "ws-message-dispatcher"
	configDynamoRegion              = "dynamodb.region"
	configDynamoUsersTableName      = "dynamodb.users-table"
	configDynamoServersTableName    = "dynamodb.servers-table"
	configDynamoChatConfigTableName = "dynamodb.chat-config-table"
	configLambdaRegion              = "lambda.region"
	configLambdaFunctionName        = "lambda.function"
	configServiceHost               = "http.host"

	envConfigDynamoRegion              = "DYNAMODB_REGION"
	envConfigDynamoUsersTableName      = "DYNAMODB_USERS_TABLE"
	envConfigDynamoServersTableName    = "DYNAMODB_SERVERS_TABLE"
	envConfigDyanmoChatConfigTableName = "DYNAMODB_CHATCONFIG_TABLE"
	envConfigLambdaRegion              = "LAMBDA_REGION"
	envConfigLambdaFunctionName        = "LAMBDA_FUNCTION"
	envConfigServiceHost               = "HTTP_HOST"

	errMissingConfiguration       = errors.New("missing configuration")
	errUnableToReadConfigFile     = errors.New("unable to read config file")
	errEmptyDynamoRegion          = errors.New("empty dynamo region")
	errEmptyDynamoUsersTable      = errors.New("missing dynamo users table configuration")
	errEmptyDynamoServersTable    = errors.New("missing dynamo servers table")
	errEmptyDynamoChatConfigTable = errors.New("missing dynamo chat config table")
)

type dynamoConfig struct {
	Region          string
	UsersTable      string
	ServersTable    string
	ChatConfigTable string
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
			"dynamo-Region":           conf.Dynamo.Region,
			"dynamo-users-table":      conf.Dynamo.UsersTable,
			"dynamo-servers-table":    conf.Dynamo.ServersTable,
			"dynamo-chatconfig-table": conf.Dynamo.ChatConfigTable,
			"lambda-Region":           conf.Lambda.Region,
			"lambda-function":         conf.Lambda.Function,
			"http-host":               conf.Service.Host,
		}).Info("config read from file")

		return conf, nil
	}

	log.Warn("unable to read config file, trying with env vars")
	err = readFromEnv(&conf)
	if err != nil {
		return Config{}, err
	}

	log.WithFields(log.Fields{
		"dynamo-Region":           conf.Dynamo.Region,
		"dynamo-users-table":      conf.Dynamo.UsersTable,
		"dynamo-servers.table":    conf.Dynamo.ServersTable,
		"dynamo-chatconfig-table": conf.Dynamo.ChatConfigTable,
		"lambda-Region":           conf.Lambda.Region,
		"lambda-function":         conf.Lambda.Function,
		"http-host":               conf.Service.Host,
	}).Info("config read from envs")

	return conf, nil
}

func readFromEnv(conf *Config) error {
	configVars := map[string]string{
		envConfigDynamoRegion:              "",
		envConfigDynamoUsersTableName:      "",
		envConfigDynamoServersTableName:    "",
		envConfigDyanmoChatConfigTableName: "",
		envConfigLambdaRegion:              "",
		envConfigLambdaFunctionName:        "",
		envConfigServiceHost:               "",
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
	conf.Dynamo.UsersTable = configVars[envConfigDynamoUsersTableName]
	conf.Dynamo.ServersTable = configVars[envConfigDynamoServersTableName]
	conf.Dynamo.ChatConfigTable = configVars[envConfigDyanmoChatConfigTableName]
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
	viper.SetDefault(configDynamoUsersTableName, defaultDynamoUsersDBTable)
	viper.SetDefault(configDynamoServersTableName, defaultDynamoServersDBTable)
	viper.SetDefault(configDynamoChatConfigTableName, defaultDynamoChatConfigTable)
	viper.SetDefault(configServiceHost, defaultHTTPHost)

	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"config_file": configFileName,
		}).Error("unable to read config file")
		return errUnableToReadConfigFile
	}

	conf.Dynamo.Region = viper.GetString(configDynamoRegion)
	conf.Dynamo.UsersTable = viper.GetString(configDynamoUsersTableName)
	conf.Dynamo.ServersTable = viper.GetString(configDynamoServersTableName)
	conf.Dynamo.ChatConfigTable = viper.GetString(configDynamoChatConfigTableName)
	conf.Lambda.Region = viper.GetString(configLambdaRegion)
	conf.Lambda.Function = viper.GetString(configLambdaFunctionName)
	conf.Service.Host = viper.GetString(configServiceHost)

	if len(conf.Lambda.Function) == 0 {
		log.Error("lambda function does not set")
		return errMissingConfiguration
	}

	return nil
}

// GetDynamoRegion gets dynamo region
func (c Config) GetDynamoRegion() (string, error) {
	if len(c.Dynamo.Region) == 0 {
		return "", errEmptyDynamoRegion
	}
	return c.Dynamo.Region, nil
}

// GetUsersTable gets dynamo users table
func (c Config) GetUsersTable() (string, error) {
	if len(c.Dynamo.UsersTable) == 0 {
		return "", errEmptyDynamoUsersTable
	}
	return c.Dynamo.UsersTable, nil
}

// GetServersTable gets dynamo servers table
func (c Config) GetServersTable() (string, error) {
	if len(c.Dynamo.ServersTable) == 0 {
		return "", errEmptyDynamoServersTable
	}
	return c.Dynamo.ServersTable, nil
}

// GetChatConfigTable gets dynamo chat config table
func (c Config) GetChatConfigTable() (string, error) {
	if len(c.Dynamo.ChatConfigTable) == 0 {
		return "", errEmptyDynamoChatConfigTable
	}
	return c.Dynamo.ChatConfigTable, nil
}
