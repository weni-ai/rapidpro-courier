package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/nyaruka/courier/v26"
	"github.com/nyaruka/courier/v26/utils"
	"github.com/pkg/errors"
)

func SendWebhooks(channel courier.Channel, r *http.Request, configWebhook interface{}, clog *courier.ChannelLog) error {
	webhook, ok := configWebhook.(map[string]interface{})
	if !ok {
		return fmt.Errorf("conversion error")
	}

	// check if url is valid
	_, err := url.ParseRequestURI(webhook["url"].(string))
	if err != nil {
		return err
	}
	u, err := url.Parse(webhook["url"].(string))
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid url %s", err)
	}

	method := webhook["method"].(string)
	if method == "" {
		method = "POST"
	}

	req, _ := http.NewRequest(method, webhook["url"].(string), r.Body)

	headers := webhook["headers"].(map[string]interface{})
	for name, value := range headers {
		req.Header.Set(name, value.(string))
	}

	trace, _, err := utils.TraceHTTP(http.DefaultClient, req, 1024)
	if trace != nil {
		clog.HTTP(trace)
	}
	if err != nil || trace == nil || trace.Response == nil || trace.Response.StatusCode/100 != 2 {
		return errors.Wrap(err, "status other than 200")
	}

	return nil
}
