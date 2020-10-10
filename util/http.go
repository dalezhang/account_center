package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"sync"

	log "gitee.com/dalezhang/account_center/logger"
)

var (
	client *Client
)

type Client struct {
	lk sync.Mutex
	wg sync.WaitGroup

	doneChan    chan struct{}
	refreshChan chan struct{}
	exitChan    chan struct{}

	Token string
}

func NewClient() *Client {

	c := &Client{
		refreshChan: make(chan struct{}),
		exitChan:    make(chan struct{}),

		Token: Config.ZeusAdminToken,
	}
	return c
}

func (c *Client) Close() {
	if c.exitChan != nil {
		close(c.exitChan)
	}
}

func (c *Client) Get(api string, params url.Values, response interface{}) error {
	return c.doRequest(http.MethodGet, api, params, nil, response)
}

func (c *Client) Post(api string, body interface{}, response interface{}) error {
	return c.doRequest(http.MethodPost, api, nil, body, response)
}

func (c *Client) Put(api string, body interface{}, response interface{}) error {
	return c.doRequest(http.MethodPut, api, nil, body, response)
}

func (c *Client) Del(api string, response interface{}) error {
	return c.doRequest(http.MethodDelete, api, nil, nil, response)
}

// TODO other errors retry
func (c *Client) doRequest(method, api string, params url.Values, bodyParams interface{}, response interface{}) error {
	var (
		err error
	)

	err = doRequest(method, api, c.Token, params, bodyParams, response)
	if err != nil {
		return err
	}

	return nil

}

func doRequest(method, api, token string, params url.Values, bodyParams interface{}, response interface{}) error {
	var (
		body io.Reader
	)
	if bodyParams != nil {
		if data, err := json.Marshal(bodyParams); err != nil {
			return err
		} else {
			body = bytes.NewBuffer(data)
		}
	}

	contentType := "application/json"
	if method == http.MethodPost {
		if body == nil {
			contentType = "application/x-www-form-urlencoded"
		}
		if len(params) > 0 {
			body = bytes.NewBufferString(params.Encode())
		}
	}

	request, err := http.NewRequest(method, api, body)
	if err != nil {
		return err
	}

	if method == http.MethodGet || method == http.MethodDelete {
		if len(params) > 0 {
			request.URL.RawQuery = params.Encode()
		}
	}
	if method != http.MethodGet {
		request.Header.Set("Content-Type", contentType)
	}
	request.Header.Set("Authorization", token)

	if os.Getenv("DEBUG") != "" {
		log.Logger.Debugf("req: %+v", request)
	}

	resp, err := http.DefaultClient.Do(request)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if os.Getenv("DEBUG") != "" {
		log.Logger.Debugf("resp: %+v", string(bodyBytes))
	}
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("req[%+v], resp[code:%d, %s]", request, resp.StatusCode, string(bodyBytes))
		log.Logger.Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	if response == nil {
		return nil
	}

	return json.Unmarshal(bodyBytes, response)
}
