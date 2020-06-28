package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	httpMaxTimeOut = 5 * time.Second
)

type chatResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,ommitempty"`
	Count   int    `json:"delivered_messages,ommitempty"`
	Elapse  string `json:"elapse,ommitpemty"`
}

func (s service) neermeChat(message interface{}) {
	var wg sync.WaitGroup
	servers := make(map[string]int, 0)

	if err := s.dbUser.GetServerConnections(servers); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to get list of servers")
		return
	}

	if len(servers) < 1 {
		log.Error("there is not configured chat-servers")
		return
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(httpMaxTimeOut))
	defer cancel()

	for ipServer, port := range servers {
		log.WithFields(log.Fields{"server": ipServer}).Info("sending request")
		wg.Add(1)
		go neermeSendMessages(ctx, ipServer, port, message, &wg)
	}

	wg.Wait()
}

func neermeSendMessages(ctx context.Context, ip string, port int, messages interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	reqJSON, err := json.Marshal(messages)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ip":    ip,
		}).Error("unable to convert messages from interface{} to []byte")
		return
	}
	reqData := strings.NewReader(string(reqJSON))

	reqString := fmt.Sprintf("http://%s:%d/publish/chat/", ip, port)
	req, err := http.NewRequest("POST", reqString, reqData)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ip":    ip,
		}).Error("unable to create new request")
		return
	}

	req.Header.Add("Content-Type", "application/json")

	req = req.WithContext(ctx)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ip":    ip,
		}).Error("unable to send request")
		return
	}

	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ip":    ip,
		}).Error("unable to read response")
		return
	}
	defer resp.Body.Close()

	response := chatResponse{}
	if err = json.Unmarshal(body, &response); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ip":    ip,
		}).Error("unable to decode response")
		return
	}

	log.WithFields(log.Fields{
		"success": response.Success,
		"count":   response.Count,
		"elapse":  response.Elapse,
		"error":   response.Error,
	}).Info("response")
}
