/*
Mailgun client in Go.
*/
package mailgun

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	API_VERSION  = 2
	API_ENDPOINT = "api.mailgun.net"
)

type Client struct {
	httpClient *http.Client
	key        string
	logger     ErrorLogger
}

func New(key string) *Client {
	l := FmtErrorLogger{}
	return &Client{httpClient: &http.Client{}, key: key, logger: l}
}

func (c *Client) SetLogger(l ErrorLogger) {
	c.logger = l
}

// make an api request
func (c *Client) api(method string, path string, fields url.Values) (body []byte, err error) {
	var req *http.Request
	url := fmt.Sprintf("https://%s/v%d%s", API_ENDPOINT, API_VERSION, path)

	if method == "POST" && fields != nil {
		req, err = http.NewRequest(method, url, strings.NewReader(fields.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	} else {
		if fields != nil {
			url += "?" + fields.Encode()
		}
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return
	}
	req.SetBasicAuth("api", c.key)
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	if rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		//err = fmt.Errorf("mailgun error: %d %s", rsp.StatusCode, body)
		msg := string(body[:])
		err = c.logger.ErrorLog("mailgun.error", rsp.StatusCode, msg)
	}
	return
}
