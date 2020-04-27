package sender

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/lambda"
	log "github.com/sirupsen/logrus"
)

type payloadLambdaRequest struct {
	Message       interface{} `json:"message"`
	ConnectionIDS []string    `json:"connection_ids"`
}

const maxRequestPerLambda = 100
const percentageGraceSingleLambda = .2

const maxRequestGraceSingleLambda = maxRequestPerLambda * percentageGraceSingleLambda

// SendMessage send messages to ws-MessageSender lambda
func (s sender) SendMessage(connections []string, msg interface{}) {

	connectionsLen := len(connections)

	connectionsPerLambda := connectionsLen / maxRequestPerLambda
	connectionsPerLambdaMod := connectionsLen % maxRequestPerLambda

	suitableForSingleLambda := connectionsPerLambda <= 1 || (connectionsLen == 1 && connectionsPerLambdaMod <= maxRequestGraceSingleLambda)

	log.WithFields(log.Fields{
		"suitableForSingleLambda": suitableForSingleLambda,
	}).Info("suitableForSingleLambda")

	if suitableForSingleLambda {
		payload := payloadLambdaRequest{
			Message:       msg,
			ConnectionIDS: connections,
		}
		s.LambdaHandler(payload)
	} else {
		for idx := 0; idx < connectionsLen; idx += maxRequestPerLambda {

			payload := payloadLambdaRequest{
				Message: msg,
			}

			sliceEnd := idx + maxRequestPerLambda

			if sliceEnd > connectionsLen {
				payload.ConnectionIDS = connections[idx:connectionsLen]
			} else {
				payload.ConnectionIDS = connections[idx:sliceEnd]
			}
			go s.LambdaHandler(payload)
		}
	}

}

func (s sender) LambdaHandler(payload payloadLambdaRequest) {

	log.WithFields(log.Fields{
		"sending": "i am in lambdahandler",
	}).Info("sending")

	payloadJSON, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("Json Marshalling error")
	}

	input := &lambda.InvokeInput{
		FunctionName: s.lambdaName,
		Payload:      payloadJSON,
	}

	result, err := s.Invoke(input)

	if err != nil {
		log.WithFields(log.Fields{
			"sending": "error sending",
		}).Error(err)
	}
	log.WithFields(log.Fields{
		"payload": string(result.Payload),
	}).Info(string(result.Payload))

}
