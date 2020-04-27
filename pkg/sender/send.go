package sender

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/lambda"
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

	payloadJSON, err := json.Marshal(payload)

	if err != nil {
		fmt.Println("Json Marshalling error")
	}

	input := &lambda.InvokeInput{
		FunctionName: s.lambdaName,
		Payload:      payloadJSON,
	}

	s.Invoke(input)

}
