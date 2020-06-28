package dynamodb

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var (
	chatTypeKeyLabel   = "label"
	chatTypeKeyValue   = "chat-type"
	chatTypeValueLabel = "value"

	errNoChatTypeConfig = errors.New("no config for chat type found")
)

func (db storage) GetChatType() (string, error) {
	projection := expression.NamesList(expression.Name(chatTypeValueLabel))
	filter := expression.Name(chatTypeKeyLabel).Equal(expression.Value(chatTypeKeyValue))

	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		return "", err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(db.chatConfigTable),
	}

	output, err := db.Scan(input)
	item := output.Items
	for _, itemConf := range item {
		if value, exists := itemConf[chatTypeValueLabel]; exists {
			return *value.S, nil
		}
	}

	return "", errNoChatTypeConfig
}
