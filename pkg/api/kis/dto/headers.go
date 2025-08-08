package dto

import (
	"auto-trader/pkg/shared/utils"
	"net/http"
)

// KISHeaders KIS API 공통 헤더 구조체
type KISHeaders struct {
	ContentType   string `json:"content-type"`
	Authorization string `json:"authorization" validate:"required,min=1,max=350"`
	AppKey        string `json:"appkey" validate:"required,min=1,max=36"`
	AppSecret     string `json:"appsecret" validate:"required,min=1,max=180"`
	TrID          string `json:"tr_id" validate:"required,min=1,max=13"`
	TrCont        string `json:"tr_cont"`
	CustType      string `json:"custtype" validate:"required,enum=P,B"`
	HashKey       string `json:"hashkey"`
}

// NewKISHeaders 새로운 KIS 헤더 생성 (개인 고객용)
func NewKISHeaders(appKey, appSecret, accessToken, trID, hashKey string) *KISHeaders {
	return &KISHeaders{
		ContentType:   "application/json; charset=utf-8",
		Authorization: "Bearer " + accessToken,
		AppKey:        appKey,
		AppSecret:     appSecret,
		TrID:          trID,
		TrCont:        "",  // 초기 조회
		CustType:      "P", // 개인
		HashKey:       hashKey,
	}
}

// Validate KISHeaders 검증
func (h *KISHeaders) Validate() error {
	return utils.ValidateStruct(h)
}

// ApplyToRequest HTTP 요청에 KIS 헤더 적용
func (h *KISHeaders) ApplyToRequest(req *http.Request) {
	req.Header.Set("Content-Type", h.ContentType)
	req.Header.Set("authorization", h.Authorization)
	req.Header.Set("appKey", h.AppKey)
	req.Header.Set("appSecret", h.AppSecret)
	req.Header.Set("tr_id", h.TrID)
	req.Header.Set("tr_cont", h.TrCont)
	req.Header.Set("custtype", h.CustType)
	if h.HashKey != "" {
		req.Header.Set("hashkey", h.HashKey)
	}
}

// TR IDs 상수 정의
const (
	// 실전투자 TR IDs
	TrIDOverseasBalanceReal = "TTTS3012R"     // 해외주식 잔고조회 (실전)
	TrIDOverseasPriceReal   = "HHDFS00000300" // 해외주식 현재가 (실전)

	// 모의투자 TR IDs
	TrIDOverseasBalanceDemo = "VTTS3012R"     // 해외주식 잔고조회 (모의)
	TrIDOverseasPriceDemo   = "HHDFS00000300" // 해외주식 현재가 (모의)
)

// NewBalanceHeaders 잔고 조회용 헤더 생성
func NewBalanceHeaders(appKey, appSecret, accessToken, hashKey string, isDemo bool) *KISHeaders {
	trID := TrIDOverseasBalanceReal
	if isDemo {
		trID = TrIDOverseasBalanceDemo
	}

	return NewKISHeaders(appKey, appSecret, accessToken, trID, hashKey)
}

// NewPriceHeaders 현재가 조회용 헤더 생성
func NewPriceHeaders(appKey, appSecret, accessToken, hashKey string, isDemo bool) *KISHeaders {
	trID := TrIDOverseasPriceReal
	if isDemo {
		trID = TrIDOverseasPriceDemo
	}

	return NewKISHeaders(appKey, appSecret, accessToken, trID, hashKey)
}
