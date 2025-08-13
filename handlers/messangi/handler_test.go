package messangi

import (
	"testing"

	"github.com/nyaruka/courier"
	. "github.com/nyaruka/courier/handlers"
	"github.com/nyaruka/courier/test"
	"github.com/nyaruka/courier/utils/clogs"
	"github.com/nyaruka/gocommon/httpx"
	"github.com/nyaruka/gocommon/urns"
)

var testChannels = []courier.Channel{
	test.NewMockChannel("8eb23e93-5ecb-45ba-b726-3b064e0c56ab", "MG", "12345", "BR", []string{urns.Phone.Prefix},
		map[string]any{
			courier.ConfigAuthToken: "test-auth-token",
		}),
}

const (
	receiveURL = "/c/mg/8eb23e93-5ecb-45ba-b726-3b064e0c56ab/receive/"
)

var testCases = []IncomingTestCase{
	{
		Label:                "Receive Valid",
		URL:                  receiveURL,
		Data:                 `{"owner":"empresa","date":"2025-05-23T14:30:00Z","processId":123,"origin":"12345","externalId":"campanha_001","callback":"https://suaapi.com/messangi/mo","connection":"SMS","id":"usuario_456","text":"Olá, quero participar","user":"+5588999999999","extraInfo":null}`,
		ExpectedRespStatus:   200,
		ExpectedBodyContains: "Message Accepted",
		ExpectedMsgText:      Sp("Olá, quero participar"),
		ExpectedURN:          "tel:+5588999999999",
		ExpectedExternalID:   "usuario_456",
	},
	{
		Label:                "Receive Missing User",
		URL:                  receiveURL,
		Data:                 `{"text":"Hello"}`,
		ExpectedRespStatus:   400,
		ExpectedBodyContains: "missing required field 'user'",
	},
	{
		Label:                "Receive Missing Text",
		URL:                  receiveURL,
		Data:                 `{"user":"+5588999999999"}`,
		ExpectedRespStatus:   400,
		ExpectedBodyContains: "missing required field 'text'",
	},
}

func TestIncoming(t *testing.T) {
	RunIncomingTestCases(t, testChannels, newHandler(), testCases)
}

func BenchmarkHandler(b *testing.B) {
	RunChannelBenchmarks(b, testChannels, newHandler(), testCases)
}

var defaultSendTestCases = []OutgoingTestCase{
	{
		Label:   "Plain Send",
		MsgText: "Simple Message ☺",
		MsgURN:  "tel:+5588999999999",
		MockResponses: map[string][]*httpx.MockResponse{
			"https://elastic.messangi.me/raven/v2/messages": {
				httpx.NewMockResponse(200, nil, []byte(`{
					"meta": {
						"timestamp": 1629174338140,
						"transactionId": "c4cdd080-fbbf-46e3-aff1-311e21f5048c"
					},
					"data": {
						"id": "abc123-def456-ghi789",
						"from": "12345",
						"text": "Simple Message ☺",
						"to": "+5588999999999",
						"type": "MT",
						"status": "ACCEPTED",
						"date": "2021-08-17T04:25:38.117592194Z"
					}
				}`)),
			},
		},
		ExpectedRequests: []ExpectedRequest{{
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer test-auth-token",
			},
			Body: `{"from":"12345","text":"Simple Message ☺","to":"5588999999999","type":"MT"}`,
		}},
		ExpectedExtIDs: []string{"abc123-def456-ghi789"},
	},
	{
		Label:   "Connection Error",
		MsgText: "Connection Error",
		MsgURN:  "tel:+5588999999999",
		MockResponses: map[string][]*httpx.MockResponse{
			"https://elastic.messangi.me/raven/v2/messages": {
				httpx.NewMockResponse(500, nil, []byte(`Internal Server Error`)),
			},
		},
		ExpectedRequests: []ExpectedRequest{{
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer test-auth-token",
			},
			Body: `{"from":"12345","text":"Connection Error","to":"5588999999999","type":"MT"}`,
		}},
		ExpectedError: courier.ErrConnectionFailed,
	},
	{
		Label:   "Invalid JSON Response",
		MsgText: "Invalid JSON",
		MsgURN:  "tel:+5588999999999",
		MockResponses: map[string][]*httpx.MockResponse{
			"https://elastic.messangi.me/raven/v2/messages": {
				httpx.NewMockResponse(200, nil, []byte(`invalid json response`)),
			},
		},
		ExpectedRequests: []ExpectedRequest{{
			Headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer test-auth-token",
			},
			Body: `{"from":"12345","text":"Invalid JSON","to":"5588999999999","type":"MT"}`,
		}},
		ExpectedLogErrors: []*clogs.LogError{
			clogs.NewLogError("response_unparseable", "", "Unable to parse response body from Messangi"),
		},
		ExpectedError: courier.ErrResponseUnparseable,
	},
}

func TestOutgoing(t *testing.T) {
	var defaultChannel = test.NewMockChannel("8eb23e93-5ecb-45ba-b726-3b064e0c56ab", "MG", "12345", "BR",
		[]string{urns.Phone.Prefix},
		map[string]any{
			courier.ConfigAuthToken: "test-auth-token",
		})
	RunOutgoingTestCases(t, defaultChannel, newHandler(), defaultSendTestCases, []string{"test-auth-token"}, nil)
}
