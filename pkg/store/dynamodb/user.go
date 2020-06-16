package dynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const (
	connectionIDLabel   = "connection_id"
	eventSubdomainLabel = "event_subdomain"
	isOrganizerLabel    = "is_organizer"
	audienceOrganizer   = "organizer"
	audienceAttendance  = "attendance"
)

func (db storage) GetUserConnections(subdomain string, audienceType string, connections *[]string) error {
	var isOrganizer bool

	switch audienceType {
	case audienceOrganizer:
		isOrganizer = true
	case audienceAttendance:
		isOrganizer = false
	default:
		audienceType = ""
	}

	var filter expression.ConditionBuilder
	if audienceType != "" {
		filter = expression.Name(eventSubdomainLabel).Equal(expression.Value(subdomain)).
			And(expression.Name(isOrganizerLabel).Equal(expression.Value(isOrganizer)))
	} else {
		filter = expression.Name(eventSubdomainLabel).Equal(expression.Value(subdomain))
	}

	projection := expression.NamesList(expression.Name(connectionIDLabel))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		return err
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(db.table),
	}

	db.ScanPages(input, func(output *dynamodb.ScanOutput, lastPage bool) bool {
		if err = appendResults(output.Items, connections); err != nil {
			return false
		}
		return true
	})

	if err != nil {
		return err
	}

	return nil
}

func appendResults(items []map[string]*dynamodb.AttributeValue, connections *[]string) error {
	conns := []struct {
		ConectionID string `json:"connection_id"`
	}{}

	if err := dynamodbattribute.UnmarshalListOfMaps(items, &conns); err != nil {
		return err
	}

	for _, conn := range conns {
		*connections = append(*connections, conn.ConectionID)
	}

	return nil
}
