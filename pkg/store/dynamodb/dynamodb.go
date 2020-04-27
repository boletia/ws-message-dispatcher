package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	awsRegion     = "us-east-1"
	dynamodbTable = "streaming-users-online"
)

type storage struct {
	*dynamodb.DynamoDB
}

// New creates new dynamodb client
func New() storage {
	sess := session.Must(
		session.NewSession(&aws.Config{
			Region: &awsRegion},
		),
	)
	dynamoDB := dynamodb.New(sess)

	return storage{
		dynamoDB,
	}
}
