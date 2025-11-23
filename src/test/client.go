package test

import (
	"avito-test/internal/handler"
	"fmt"
	"log"

	"gopkg.in/resty.v1"
)

type ShopClient struct {
	httpClient *resty.Client
	baseUrl    string
	authToken  *string
}

const (
	AUTH_ROUTE     = "/api/auth"
	BUY_ITEM_ROUTE = "/api/buy"
	INFO_ROUTE     = "/api/info"
	SEND_ROUTE     = "/api/sendCoin"
)

func NewShopClient(baseURL string) *ShopClient {
	return &ShopClient{
		httpClient: resty.New(),
		baseUrl:    baseURL,
		authToken:  nil,
	}
}

func (c *ShopClient) GetInfo() (int, handler.InfoResponse) {
	var info handler.InfoResponse

	code := c.callGet(c.baseUrl+INFO_ROUTE, &info)

	return code, info
}

func (c *ShopClient) BuyItem(item string) int {
	var empty any

	code := c.callGet(fmt.Sprintf("%s%s/%s", c.baseUrl, BUY_ITEM_ROUTE, item), &empty)

	return code
}

func (c *ShopClient) Auth(username string, password string) (int, handler.AuthResponse) {
	var req = handler.AuthRequest{
		Username: username,
		Password: password,
	}
	var response handler.AuthResponse

	code := c.callPost(c.baseUrl+AUTH_ROUTE, &req, &response)

	return code, response
}

func (c *ShopClient) Send(to_username string, amount uint) int {
	var req = handler.SendCoinRequest{
		ToUser: to_username,
		Amount: amount,
	}
	var empty any

	code := c.callPost(c.baseUrl+SEND_ROUTE, &req, &empty)

	return code
}

func (c *ShopClient) SetAuthToken(token string) {
	c.authToken = &token
}

func (c *ShopClient) createRestyRequest(res any) *resty.Request {
	request := c.httpClient.R().
		SetHeader("Accept", "application/json").
		SetResult(res)

	if c.authToken != nil {
		request.SetAuthToken(*c.authToken)
	}

	return request
}

func (c *ShopClient) callGet(url string, res any) int {
	resp, err := c.createRestyRequest(res).Get(url)

	if err != nil {
		log.Printf("Error occurred: %v", err)
	}

	return resp.StatusCode()
}

func (c *ShopClient) callPost(url string, body any, res any) int {
	resp, err := c.createRestyRequest(res).SetBody(body).Post(url)

	if err != nil {
		log.Printf("Error occurred: %v", err)
	}

	return resp.StatusCode()
}
