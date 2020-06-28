package service

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

const (
	apiGatewayChat = "api-gateway"
	neermeChat     = "neerme-v2"
)

type incomeMessage struct {
	EventSubdomain string      `json:"event_subdomain"`
	AudienceType   string      `json:"audience_type"`
	Message        interface{} `json:"message"`
}

type response struct {
	Success bool `json:"success"`
}

// TakeIn receives new messages from Ws-message-connector
func (s service) TakeIn(c echo.Context) error {
	incomeMsg := incomeMessage{}

	if err := c.Bind(&incomeMsg); err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"incomeMsg": fmt.Sprintf("%#v", incomeMsg),
		}).Error("unable to decode request")
		return c.JSON(http.StatusBadRequest, response{false})
	}

	if len(incomeMsg.EventSubdomain) == 0 {
		log.Error("empty event subdomain")
		return c.JSON(http.StatusBadRequest, response{false})
	}

	log.WithFields(log.Fields{"event_subdomain": incomeMsg.EventSubdomain}).Info("request decoded")
	go s.dispatchMessage(incomeMsg)

	return c.JSON(http.StatusOK, response{true})
}

func (s service) dispatchMessage(msg incomeMessage) {
	var connections []string

	chatType, err := s.dbUser.GetChatType()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to get chat type")
		return
	}

	switch chatType {
	case apiGatewayChat:
		log.WithFields(log.Fields{"chat-type": apiGatewayChat}).Info("sending messages")
		if err := s.dbUser.GetUserConnections(msg.EventSubdomain, msg.AudienceType, &connections); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("unable to get user connections")
			return
		}

		log.WithFields(log.Fields{
			"connections": connections,
		}).Info("connections")

		s.sender.SendMessage(connections, msg.Message)

	case neermeChat:
		log.WithFields(log.Fields{"chat-type": neermeChat}).Info("sending messages")
		s.neermeChat(msg)

	default:
		log.WithFields(log.Fields{
			"type": chatType,
		}).Error("unknow chat type")
	}
}
