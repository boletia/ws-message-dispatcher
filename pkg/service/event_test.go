package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestTakeIn(t *testing.T) {
	testCases := []struct {
		testName               string
		requestPost            string
		expectedResponse       string
		expectedHTTPStatusCode int
	}{
		{
			testName:               "DefaultRegularCase",
			requestPost:            `{ "event_subdomain":"el-show-de-producto-online", "audience_type":"attendance", "message": { "active": false, "answers": [ { "id": "ANSWER-ID-0", "option_label": "indeed", "total": 1 }, { "id": "ANSWER-ID-1", "option_label": "indeednt", "total": 0 } ] } }`,
			expectedResponse:       "{\"success\":true}\n",
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			testName:               "ApiGatewayRegularCase",
			requestPost:            `{ "gateway_type": "api-gateway", "event_subdomain":"el-show-de-producto-online", "audience_type":"attendance", "message": { "active": false, "answers": [ { "id": "ANSWER-ID-0", "option_label": "indeed", "total": 1 }, { "id": "ANSWER-ID-1", "option_label": "indeednt", "total": 0 } ] } }`,
			expectedResponse:       "{\"success\":true}\n",
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			testName:               "NeermeV2RegularCase",
			requestPost:            `{ "gateway_type": "chat-server-v2", "event_subdomain":"el-show-de-producto-online", "audience_type":"attendance", "message": { "active": false, "answers": [ { "id": "ANSWER-ID-0", "option_label": "indeed", "total": 1 }, { "id": "ANSWER-ID-1", "option_label": "indeednt", "total": 0 } ] } }`,
			expectedResponse:       "{\"success\":true}\n",
			expectedHTTPStatusCode: http.StatusOK,
		},
		{
			testName:               "NoSubdomainCase",
			requestPost:            `{ "event_subdomain":"", "audience_type":"attendance", "message": { "active": false, "answers": [ { "id": "ANSWER-ID-0", "option_label": "indeed", "total": 1 }, { "id": "ANSWER-ID-1", "option_label": "indeednt", "total": 0 } ] } }`,
			expectedResponse:       "{\"success\":false}\n",
			expectedHTTPStatusCode: http.StatusBadRequest,
		},
		{
			testName:               "ErrorRequestDecodeCase",
			requestPost:            `{ event_subdomain:"el-show-de-producto-online", "audience_type":"attendance", "message": { "active": false, "answers": [ { "id": "ANSWER-ID-0", "option_label": "indeed", "total": 1 }, { "id": "ANSWER-ID-1", "option_label": "indeednt", "total": 0 } ] } }`,
			expectedResponse:       "{\"success\":false}\n",
			expectedHTTPStatusCode: http.StatusBadRequest,
		},
	}

	for _, c := range testCases {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/v1/dispatcher/event/", strings.NewReader(c.requestPost))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		context := e.NewContext(req, rec)

		userStorage := connGetter{}
		mSender := msgSender{}
		srv := New(userStorage, mSender)

		t.Run(c.testName, func(t *testing.T) {
			if assert.NoError(t, srv.TakeIn(context)) {
				assert.Equal(t, c.expectedHTTPStatusCode, rec.Code)
				assert.Equal(t, c.expectedResponse, rec.Body.String())
			}
		})
	}
}

type connGetter struct{}

func (cg connGetter) GetUserConnections(eventSubdomain string, audienceType string, connections *[]string) error {
	return nil
}

func (cg connGetter) GetServerConnections(servers map[string]int) error {
	return nil
}

func (cg connGetter) GetChatType() (string, error) {
	return "", nil
}

type msgSender struct{}

func (ms msgSender) SendMessage(connections []string, msg interface{}) {

}
