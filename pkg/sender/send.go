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
		"ConnectionIDS Length":    connectionsLen,
	}).Info("SendMessage")

	if suitableForSingleLambda {
		payload := payloadLambdaRequest{
			Message:       msg,
			ConnectionIDS: connections,
		}
		log.WithFields(log.Fields{
			"sending_to_n_connections": suitableForSingleLambda,
		}).Info("SendMessage")
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

			log.WithFields(log.Fields{
				"sending_to_n_connections": len(payload.ConnectionIDS),
			}).Info("SendMessage")

			go s.LambdaHandler(payload)
		}
	}

}

func (s sender) LambdaHandler(payload payloadLambdaRequest) {

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
			"error": err,
		}).Error("LambdaHandler")
	}
	log.WithFields(log.Fields{
		"result_lambda": result,
	}).Error("LambdaHandler")
}
