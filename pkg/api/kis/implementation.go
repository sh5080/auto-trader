package kis

import (
	"fmt"
	"time"

	"auto-trader/pkg/domain/portfolio"

	"github.com/shopspring/decimal"
)

// KISDataSource KIS API를 사용하는 외부 데이터 소스 구현체
type KISDataSource struct {
	client  *Client
	adapter *Adapter
}

// NewKISDataSource 새로운 KIS 데이터 소스 생성
func NewKISDataSource(appKey, appSecret, baseURL string, isDemo bool) *KISDataSource {
	return &KISDataSource{
		client:  NewClient(appKey, appSecret, baseURL, isDemo),
		adapter: NewAdapter(),
	}
}

// SetAccessToken Access Token 설정
func (k *KISDataSource) SetAccessToken(token string) {
	k.client.SetAccessToken(token)
}

// GetBalance 잔고 조회
func (k *KISDataSource) GetBalance(accountNo string) ([]*portfolio.Position, error) {
	// KIS API 호출
	balanceResp, err := k.client.GetBalance(accountNo)
	if err != nil {
		return nil, fmt.Errorf("KIS API 잔고 조회 실패: %w", err)
	}

	// 응답을 도메인 모델로 변환
	var positions []*portfolio.Position
	for _, kisBalance := range balanceResp.Output1 {
		position, err := k.adapter.ConvertBalanceToPosition(kisBalance)
		if err != nil {
			return nil, fmt.Errorf("포지션 변환 실패: %w", err)
		}
		positions = append(positions, position)
	}

	return positions, nil
}

// GetCurrentPrice 현재가 조회
func (k *KISDataSource) GetCurrentPrice(symbol string) (*portfolio.StockPrice, error) {
	// KIS API 호출
	priceResp, err := k.client.GetCurrentPrice(symbol)
	if err != nil {
		return nil, fmt.Errorf("KIS API 현재가 조회 실패: %w", err)
	}

	// 응답을 도메인 모델로 변환
	stockPrice, err := k.adapter.ConvertPriceToStockPrice(priceResp.Output)
	if err != nil {
		return nil, fmt.Errorf("주식 가격 변환 실패: %w", err)
	}

	return stockPrice, nil
}

// GetCurrentPrices 여러 종목 현재가 조회
func (k *KISDataSource) GetCurrentPrices(symbols []string) ([]*portfolio.StockPrice, error) {
	var stockPrices []*portfolio.StockPrice

	for _, symbol := range symbols {
		stockPrice, err := k.GetCurrentPrice(symbol)
		if err != nil {
			return nil, fmt.Errorf("종목 %s 현재가 조회 실패: %w", symbol, err)
		}
		stockPrices = append(stockPrices, stockPrice)
	}

	return stockPrices, nil
}

// GetPortfolioSummary 포트폴리오 요약 조회
func (k *KISDataSource) GetPortfolioSummary(userID, accountNo string) (*portfolio.PortfolioSummary, error) {
	// 잔고 조회
	positions, err := k.GetBalance(accountNo)
	if err != nil {
		return nil, fmt.Errorf("잔고 조회 실패: %w", err)
	}

	if len(positions) == 0 {
		return &portfolio.PortfolioSummary{
			TotalPositions:  0,
			TotalValue:      decimal.Zero,
			TotalProfit:     decimal.Zero,
			TotalProfitRate: decimal.Zero,
			DailyProfit:     decimal.Zero,
			DailyProfitRate: decimal.Zero,
			Positions:       []portfolio.Position{},
			LastUpdated:     time.Now(),
			DataFreshness:   "REALTIME",
		}, nil
	}

	// 요약 정보 계산
	totalValue := decimal.Zero
	totalProfit := decimal.Zero
	var topGainers, topLosers []portfolio.Position

	for _, pos := range positions {
		totalValue = totalValue.Add(pos.TotalValue)
		totalProfit = totalProfit.Add(pos.TotalProfit)

		// 상위 수익/손실 종목 분류
		if pos.ProfitRate.GreaterThan(decimal.Zero) {
			topGainers = append(topGainers, *pos)
		} else {
			topLosers = append(topLosers, *pos)
		}
	}

	// 수익률 계산
	totalProfitRate := decimal.Zero
	if totalValue.GreaterThan(decimal.Zero) {
		totalProfitRate = totalProfit.Div(totalValue).Mul(decimal.NewFromInt(100))
	}

	// 일일 수익은 별도 계산 필요 (현재는 0으로 설정)
	dailyProfit := decimal.Zero
	dailyProfitRate := decimal.Zero

	return &portfolio.PortfolioSummary{
		TotalPositions:  len(positions),
		TotalValue:      totalValue,
		TotalProfit:     totalProfit,
		TotalProfitRate: totalProfitRate,
		DailyProfit:     dailyProfit,
		DailyProfitRate: dailyProfitRate,
		TopGainers:      topGainers,
		TopLosers:       topLosers,
		Positions:       convertToPositions(positions),
		LastUpdated:     time.Now(),
		DataFreshness:   "REALTIME",
	}, nil
}

// convertToPositions 포인터 슬라이스를 값 슬라이스로 변환
func convertToPositions(positions []*portfolio.Position) []portfolio.Position {
	result := make([]portfolio.Position, len(positions))
	for i, pos := range positions {
		result[i] = *pos
	}
	return result
}
