package service

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
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

	// Get user connection's list
	if err := s.dbUser.GetUserConnections(msg.EventSubdomain, msg.AudienceType, &connections); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to get user connections")
		return
	}

	log.WithFields(log.Fields{
		"connections": connections,
	}).Info("connections")

	// Calculate batch send message
	// call MessageSender lambda process

}
