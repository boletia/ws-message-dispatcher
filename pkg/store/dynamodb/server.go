package dynamodb

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

var (
	serversIDLabel   = "ip"
	serversPortLabel = "port"
)

func (db storage) GetServerConnections(servers map[string]int) error {

	projection := expression.NamesList(expression.Name(serversIDLabel), expression.Name(serversPortLabel))
	expr, err := expression.NewBuilder().WithProjection(projection).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(db.serversTable),
	}

	output := &dynamodb.ScanOutput{}
	for output, err = db.Scan(input); len(output.LastEvaluatedKey) != 0 && err == nil; output, err = db.Scan(input) {
		if err = appendChatServers(output.Items, servers); err != nil {
			return err
		}

		input = &dynamodb.ScanInput{
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			TableName:                 aws.String(db.serversTable),
		}
		input.ExclusiveStartKey = output.LastEvaluatedKey
	}

	if *output.Count > 0 {
		if err = appendChatServers(output.Items, servers); err != nil {
			return err
		}
	}

	return nil
}

func appendChatServers(items []map[string]*dynamodb.AttributeValue, servers map[string]int) error {
	for _, item := range items {
		if attrIP, exists := item[serversIDLabel]; exists {
			if attrPort, exists := item[serversPortLabel]; exists {
				ip := *attrIP.S
				port, err := strconv.Atoi(*attrPort.N)
				if err != nil {
					return err
				}

				if len(ip) > 0 && port > 0 {
					servers[ip] = port
				}
			}
		}
	}
	return nil
}
