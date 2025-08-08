package dto

import "auto-trader/pkg/shared/utils"

// BalanceRequest 해외주식 잔고 조회 요청
type BalanceRequest struct {
	CANO           string `json:"CANO" validate:"required,min=1,max=20"`        // 종합계좌번호
	ACNT_PRDT_CD   string `json:"ACNT_PRDT_CD" validate:"required,min=1,max=2"` // 계좌상품코드
	OVRS_EXCG_CD   string `json:"OVRS_EXCG_CD"`                                 // 해외거래소코드
	TR_CRCY_CD     string `json:"TR_CRCY_CD"`                                   // 거래통화코드
	CTX_AREA_FK200 string `json:"CTX_AREA_FK200"`                               // 연속조회검색조건
	CTX_AREA_NK200 string `json:"CTX_AREA_NK200"`                               // 연속조회키
}

// NewBalanceRequest 새로운 잔고 조회 요청 생성
func NewBalanceRequest(accountNo string) *BalanceRequest {
	return &BalanceRequest{
		CANO:           accountNo,
		ACNT_PRDT_CD:   "01", // 기본값
		OVRS_EXCG_CD:   "",   // 전체
		TR_CRCY_CD:     "",   // 전체
		CTX_AREA_FK200: "",
		CTX_AREA_NK200: "",
	}
}

// Validate BalanceRequest 검증
func (r *BalanceRequest) Validate() error {
	return utils.ValidateStruct(r)
}

// PriceRequest 현재가 조회 요청
type PriceRequest struct {
	AUTH string `json:"AUTH"`                                  // 인증 정보
	EXCD string `json:"EXCD" validate:"required,min=1,max=10"` // 거래소코드
	SYMB string `json:"SYMB" validate:"required,min=1,max=20"` // 종목 심볼
}

// NewPriceRequest 새로운 현재가 조회 요청 생성
func NewPriceRequest(symbol string) *PriceRequest {
	return &PriceRequest{
		AUTH: "",
		EXCD: "NAS", // NASDAQ 기본값
		SYMB: symbol,
	}
}

// Validate PriceRequest 검증
func (r *PriceRequest) Validate() error {
	return utils.ValidateStruct(r)
}

// OrderRequest 주문 요청 (향후 확장용)
type OrderRequest struct {
	Symbol    string `json:"symbol" validate:"required,min=1,max=20"`     // 종목 심볼
	Quantity  int    `json:"quantity" validate:"required"`                // 수량
	Price     string `json:"price" validate:"required,min=1,max=20"`      // 가격
	OrderType string `json:"orderType" validate:"required,enum=BUY,SELL"` // 주문 타입 (BUY, SELL)
}

// NewOrderRequest 새로운 주문 요청 생성
func NewOrderRequest(symbol string, quantity int, price string, orderType string) *OrderRequest {
	return &OrderRequest{
		Symbol:    symbol,
		Quantity:  quantity,
		Price:     price,
		OrderType: orderType,
	}
}

// Validate OrderRequest 검증
func (r *OrderRequest) Validate() error {
	return utils.ValidateStruct(r)
}
