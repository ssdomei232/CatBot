package napcathttp

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(token string, url string) *Client {
	restyClient := resty.New()
	restyClient.SetBaseURL(url)
	restyClient.SetHeader("Authorization", "Bearer "+token)
	restyClient.SetHeader("User-Agent", "napcat-http-go-client/1.0")
	restyClient.SetHeader("Accept", "application/json")

	return &Client{
		client: restyClient,
	}
}

// DoRequest 发起一次http请求
func (c *Client) DoRequest(method, endpoint string, reqData any, respData any) error {
	req := c.client.R()

	if reqData != nil {
		req.SetBody(reqData)
	}

	if respData != nil {
		req.SetResult(respData)
	}

	_, err := req.Execute(method, endpoint)

	return err
}
