package sender

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type sender struct {
	*lambda.Lambda
	lambdaName *string
}

// New Creates new sender instance
func New(region, funcName string) sender {

	sess := session.New(&aws.Config{
		Region: &region,
	})
	Lambda := lambda.New(sess)

	return sender{
		Lambda,
		aws.String(funcName),
	}
}
