package api

import (
	"context"
	"fmt"
	"time"

	"auto-trader/pkg/shared/config"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type Client struct {
	httpClient *resty.Client
	config     *config.Config
}

type PriceData struct {
	Symbol    string          `json:"symbol"`
	Price     decimal.Decimal `json:"price"`
	Volume    decimal.Decimal `json:"volume"`
	Timestamp time.Time       `json:"timestamp"`
}

type OrderRequest struct {
	Symbol   string          `json:"symbol"`
	Side     string          `json:"side"` // "buy" or "sell"
	Quantity decimal.Decimal `json:"quantity"`
	Price    decimal.Decimal `json:"price,omitempty"`
	Type     string          `json:"type"` // "market" or "limit"
}

type OrderResponse struct {
	ID       string          `json:"id"`
	Symbol   string          `json:"symbol"`
	Side     string          `json:"side"`
	Quantity decimal.Decimal `json:"quantity"`
	Price    decimal.Decimal `json:"price"`
	Status   string          `json:"status"`
	Type     string          `json:"type"`
}

func NewClient(cfg *config.Config) *Client {
	client := resty.New()
	client.SetTimeout(cfg.Trading.OrderTimeout)
	client.SetBaseURL(cfg.API.BaseURL)
	client.SetHeader("Authorization", "Bearer "+cfg.API.APIKey)
	client.SetHeader("Content-Type", "application/json")

	return &Client{
		httpClient: client,
		config:     cfg,
	}
}

func (c *Client) GetCurrentPrice(symbol string) (*PriceData, error) {
	resp, err := c.httpClient.R().
		SetResult(&PriceData{}).
		Get(fmt.Sprintf("/price/%s", symbol))

	if err != nil {
		return nil, fmt.Errorf("가격 조회 실패: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API 오류: %s", resp.Status())
	}

	priceData := resp.Result().(*PriceData)
	logrus.Infof("가격 조회: %s = %s", symbol, priceData.Price.String())

	return priceData, nil
}

func (c *Client) PlaceOrder(ctx context.Context, req *OrderRequest) (*OrderResponse, error) {
	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&OrderResponse{}).
		Post("/orders")

	if err != nil {
		return nil, fmt.Errorf("주문 실패: %w", err)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
		return nil, fmt.Errorf("주문 API 오류: %s", resp.Status())
	}

	orderResp := resp.Result().(*OrderResponse)
	logrus.Infof("주문 실행: %s %s %s @ %s", orderResp.Side, orderResp.Quantity.String(), orderResp.Symbol, orderResp.Price.String())

	return orderResp, nil
}

func (c *Client) GetOrderStatus(orderID string) (*OrderResponse, error) {
	resp, err := c.httpClient.R().
		SetResult(&OrderResponse{}).
		Get(fmt.Sprintf("/orders/%s", orderID))

	if err != nil {
		return nil, fmt.Errorf("주문 상태 조회 실패: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("주문 상태 API 오류: %s", resp.Status())
	}

	return resp.Result().(*OrderResponse), nil
}

func (c *Client) CancelOrder(orderID string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/orders/%s", orderID))

	if err != nil {
		return fmt.Errorf("주문 취소 실패: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("주문 취소 API 오류: %s", resp.Status())
	}

	logrus.Infof("주문 취소: %s", orderID)
	return nil
}
