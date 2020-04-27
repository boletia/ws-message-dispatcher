package sender

type payloadLambdaRequest struct {
	Message       interface{} `json:"message"`
	ConnectionIDS []string    `json:"connection_ids"`
}

// SendMessage send messages to ws-MessageSender lambda
func (s sender) SendMessage(connections []string, msg interface{}) {
	payload := payloadLambdaRequest{
		Message: msg,
	}

	for i := 0; i < len(connections); i++ {
		payload.ConnectionIDS = append(payload.ConnectionIDS, connections[i])
		if i%50 == 0 {
			// invoke lambda with  a go routing, using wg
		}
	}
}
