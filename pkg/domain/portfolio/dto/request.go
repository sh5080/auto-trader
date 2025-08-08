package dto

// 쿼리 파라미터 DTO들

// GetPortfolioQuery 포트폴리오 조회 쿼리
type GetPortfolioQuery struct {
	ForceRefresh bool `query:"forceRefresh"`
}

// GetPortfolioSummaryQuery 포트폴리오 요약 조회 쿼리
type GetPortfolioSummaryQuery struct {
	ForceRefresh bool `query:"forceRefresh"`
}

// GetPositionsQuery 보유 주식 목록 조회 쿼리
type GetPositionsQuery struct {
	Symbol       string `query:"symbol"`
	ForceRefresh bool   `query:"forceRefresh"`
}

// GetCurrentPricesQuery 여러 종목 현재가 조회 쿼리
type GetCurrentPricesQuery struct {
	Symbols      string `query:"symbols" validate:"required,min=1"`
	ForceRefresh bool   `query:"forceRefresh"`
}

// GetTradeHistoryQuery 거래 내역 조회 쿼리
// 날짜는 컨트롤러에서 파싱하여 Service 요청 DTO로 전달
type GetTradeHistoryQuery struct {
	Symbol string `query:"symbol"`
	Limit  int    `query:"limit"`
	Offset int    `query:"offset"`
}

// 경로 파라미터 DTO들

// SymbolPath 특정 종목 심볼 경로 파라미터
type SymbolPath struct {
	Symbol string `validate:"required,min=1"`
}
