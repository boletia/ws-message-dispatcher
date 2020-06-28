package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

// ConfigSetter set config for dynamo package
type ConfigSetter interface {
	GetDynamoRegion() (string, error)
	GetUsersTable() (string, error)
	GetServersTable() (string, error)
	GetChatConfigTable() (string, error)
}

type storage struct {
	*dynamodb.DynamoDB
	usersTable      string
	serversTable    string
	chatConfigTable string
}

// New creates new dynamodb client
func New(setter ConfigSetter) storage {
	var region, usersTable, serversTable, chatConfigTable string
	var err error

	if region, err = setter.GetDynamoRegion(); err != nil {
		log.Fatalf("unable to read dynamo region:%s", err)
	}

	if usersTable, err = setter.GetUsersTable(); err != nil {
		log.Fatalf("unable to read dynamo users table:%s", err)
	}

	if serversTable, err = setter.GetServersTable(); err != nil {
		log.Fatalf("unable to read dynamo servers table:%s", err)
	}

	if chatConfigTable, err = setter.GetChatConfigTable(); err != nil {
		log.Fatalf("unable to read dynamo chat config table:%s", err)
	}

	sess := session.Must(
		session.NewSession(&aws.Config{
			Region: &region},
		),
	)
	dynamoDB := dynamodb.New(sess)

	return storage{
		dynamoDB,
		usersTable,
		serversTable,
		chatConfigTable,
	}
}
