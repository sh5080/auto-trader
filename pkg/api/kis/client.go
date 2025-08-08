package kis

import (
	"auto-trader/pkg/api/kis/dto"
	"auto-trader/pkg/shared/utils"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Client 한국투자증권 API 클라이언트
type Client struct {
	AppKey      string
	AppSecret   string
	AccessToken string
	BaseURL     string
	IsDemo      bool
	HTTPClient  *http.Client
}

// NewClient 새로운 KIS API 클라이언트 생성
func NewClient(appKey, appSecret, baseURL string, isDemo bool) *Client {
	return &Client{
		AppKey:      appKey,
		AppSecret:   appSecret,
		AccessToken: "",
		BaseURL:     baseURL,
		IsDemo:      isDemo,
		HTTPClient: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       30 * time.Second,
		},
	}
}

// generateHashkey Hashkey 생성 (한국투자증권 API 요구사항)
func (c *Client) generateHashkey(body string) (string, error) {
	h := hmac.New(sha256.New, []byte(c.AppSecret))
	h.Write([]byte(body))
	hash := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash), nil
}

// GetBalance 해외주식 잔고 조회
func (c *Client) GetBalance(accountNo string) (*KISBalanceResponse, error) {
	// API 엔드포인트
	url := fmt.Sprintf("%s/uapi/overseas-stock/v1/trading/inquire-balance", c.BaseURL)
	logrus.Infof("GetBalance URL: %s", url)
	// 요청 바디
	requestBody := dto.NewBalanceRequest(accountNo)

	// DTO 검증
	if err := requestBody.Validate(); err != nil {
		return nil, utils.WrapValidationError(err, "요청 검증 실패")
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("요청 바디 마샬링 실패: %w", err)
	}

	// Hashkey 생성
	hashkey, err := c.generateHashkey(string(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("hashkey 생성 실패: %w", err)
	}

	// HTTP 요청 생성
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("http 요청 생성 실패: %w", err)
	}

	// 공통 헤더 설정
	headers := dto.NewBalanceHeaders(c.AppKey, c.AppSecret, c.AccessToken, hashkey, c.IsDemo)

	// 헤더 검증
	if err := headers.Validate(); err != nil {
		return nil, utils.WrapValidationError(err, "헤더 검증 실패")
	}

	headers.ApplyToRequest(req)

	// 요청 실행
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API 요청 실패: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// 응답 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %w", err)
	}

	// 응답 파싱
	var balanceResp KISBalanceResponse
	if err := json.Unmarshal(body, &balanceResp); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}
	logrus.Infof("GetBalance Response: %+v", balanceResp)
	// 에러 체크
	if balanceResp.RtCd != "0" {
		return nil, fmt.Errorf("API 오류: %s - %s", balanceResp.MsgCd, balanceResp.Msg1)
	}

	return &balanceResp, nil
}

// GetCurrentPrice 해외주식 현재가 조회
func (c *Client) GetCurrentPrice(symbol string) (*KISPriceResponse, error) {
	// API 엔드포인트
	url := fmt.Sprintf("%s/uapi/overseas-price/v1/quotations/price", c.BaseURL)

	// 요청 바디
	requestBody := dto.NewPriceRequest(symbol)

	// DTO 검증
	if err := requestBody.Validate(); err != nil {
		return nil, utils.WrapValidationError(err, "요청 검증 실패")
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("요청 바디 마샬링 실패: %w", err)
	}

	// Hashkey 생성
	hashkey, err := c.generateHashkey(string(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("hashkey 생성 실패: %w", err)
	}

	// HTTP 요청 생성
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("http 요청 생성 실패: %w", err)
	}

	// 공통 헤더 설정
	headers := dto.NewPriceHeaders(c.AppKey, c.AppSecret, c.AccessToken, hashkey, c.IsDemo)

	// 헤더 검증
	if err := headers.Validate(); err != nil {
		return nil, utils.WrapValidationError(err, "헤더 검증 실패")
	}

	headers.ApplyToRequest(req)

	// 요청 실행
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API 요청 실패: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// 응답 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %w", err)
	}

	// 응답 파싱
	var priceResp KISPriceResponse
	if err := json.Unmarshal(body, &priceResp); err != nil {
		return nil, fmt.Errorf("응답 파싱 실패: %w", err)
	}

	// 에러 체크
	if priceResp.RtCd != "0" {
		return nil, fmt.Errorf("API 오류: %s - %s", priceResp.MsgCd, priceResp.Msg1)
	}

	return &priceResp, nil
}

// SetAccessToken Access Token 설정
func (c *Client) SetAccessToken(token string) {
	c.AccessToken = token
}
