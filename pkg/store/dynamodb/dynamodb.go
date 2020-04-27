package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type storage struct {
	*dynamodb.DynamoDB
	table string
}

// New creates new dynamodb client
func New(region, table string) storage {
	sess := session.Must(
		session.NewSession(&aws.Config{
			Region: &region},
		),
	)
	dynamoDB := dynamodb.New(sess)

	return storage{
		dynamoDB,
		table,
	}
}
