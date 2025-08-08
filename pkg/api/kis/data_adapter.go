package kis

import (
	"auto-trader/pkg/domain/portfolio"
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

// DataAdapter KIS API를 portfolio 도메인의 ExternalAPI 인터페이스에 맞게 어댑터
type DataAdapter struct {
	client  *Client
	adapter *Adapter
}

// NewDataAdapter 새로운 데이터 어댑터 생성
func NewDataAdapter(appKey, appSecret, baseURL string, isDemo bool) *DataAdapter {
	return &DataAdapter{
		client:  NewClient(appKey, appSecret, baseURL, isDemo),
		adapter: NewAdapter(),
	}
}

// SetAccessToken Access Token 설정
func (d *DataAdapter) SetAccessToken(token string) {
	d.client.SetAccessToken(token)
}

// GetCurrentPrice 현재가 조회
func (d *DataAdapter) GetCurrentPrice(symbol string) (*portfolio.StockPrice, error) {
	// KIS API 호출
	priceResp, err := d.client.GetCurrentPrice(symbol)
	if err != nil {
		return nil, fmt.Errorf("KIS API 현재가 조회 실패: %w", err)
	}

	// 응답을 portfolio 도메인 모델로 변환
	stockPrice, err := d.convertToDataStockPrice(priceResp.Output)
	if err != nil {
		return nil, fmt.Errorf("주식 가격 변환 실패: %w", err)
	}

	return stockPrice, nil
}

// GetCurrentPrices 여러 종목 현재가 조회
func (d *DataAdapter) GetCurrentPrices(symbols []string) ([]portfolio.StockPrice, error) {
	var stockPrices []portfolio.StockPrice

	for _, symbol := range symbols {
		stockPrice, err := d.GetCurrentPrice(symbol)
		if err != nil {
			return nil, fmt.Errorf("종목 %s 현재가 조회 실패: %w", symbol, err)
		}
		stockPrices = append(stockPrices, *stockPrice)
	}

	return stockPrices, nil
}

// GetUserPositions 사용자 보유 주식 조회
func (d *DataAdapter) GetUserPositions(userID string) ([]portfolio.Position, error) {
	// TODO: 사용자별 계좌번호 매핑 필요
	accountNo := "계좌번호" // 실제로는 사용자별 계좌번호를 가져와야 함

	// KIS API 호출
	balanceResp, err := d.client.GetBalance(accountNo)
	if err != nil {
		return nil, fmt.Errorf("KIS API 잔고 조회 실패: %w", err)
	}

	// 응답을 portfolio 도메인 모델로 변환
	var positions []portfolio.Position
	for _, kisBalance := range balanceResp.Output1 {
		position, err := d.convertToDataPosition(kisBalance, userID)
		if err != nil {
			return nil, fmt.Errorf("포지션 변환 실패: %w", err)
		}
		positions = append(positions, *position)
	}

	return positions, nil
}

// GetUserPortfolio 사용자 포트폴리오 조회
func (d *DataAdapter) GetUserPortfolio(userID string) (*portfolio.Portfolio, error) {
	positions, err := d.GetUserPositions(userID)
	if err != nil {
		return nil, err
	}

	// 포트폴리오 계산
	totalValue := decimal.Zero
	totalProfit := decimal.Zero
	for _, pos := range positions {
		totalValue = totalValue.Add(pos.TotalValue)
		totalProfit = totalProfit.Add(pos.TotalProfit)
	}

	profitRate := decimal.Zero
	if totalValue.GreaterThan(decimal.Zero) {
		profitRate = totalProfit.Div(totalValue).Mul(decimal.NewFromInt(100))
	}

	return &portfolio.Portfolio{
		UserID:      userID,
		TotalValue:  totalValue,
		TotalProfit: totalProfit,
		ProfitRate:  profitRate,
		UpdatedAt:   time.Now(),
	}, nil
}

// GetCompanyInfo 회사 정보 조회 (현재는 미구현)
func (d *DataAdapter) GetCompanyInfo(symbol string) (*portfolio.CompanyInfo, error) {
	// TODO: KIS API에서 회사 정보 조회 구현
	return nil, fmt.Errorf("회사 정보 조회는 아직 구현되지 않았습니다")
}

// GetChartData 차트 데이터 조회 (현재는 미구현)
func (d *DataAdapter) GetChartData(symbol string, period string, startDate, endDate time.Time) ([]portfolio.ChartData, error) {
	// TODO: KIS API에서 차트 데이터 조회 구현
	return nil, fmt.Errorf("차트 데이터 조회는 아직 구현되지 않았습니다")
}

// GetTradeHistory 거래 내역 조회 (현재는 미구현)
func (d *DataAdapter) GetTradeHistory(userID string, symbol string, startDate, endDate time.Time) ([]portfolio.TradeHistory, error) {
	// TODO: KIS API에서 거래 내역 조회 구현
	return nil, fmt.Errorf("거래 내역 조회는 아직 구현되지 않았습니다")
}

// convertToDataStockPrice KIS 현재가 데이터를 portfolio.StockPrice로 변환
func (d *DataAdapter) convertToDataStockPrice(kisPrice KISPriceOutput) (*portfolio.StockPrice, error) {
	price, err := decimal.NewFromString(kisPrice.Last)
	if err != nil {
		return nil, fmt.Errorf("현재가 파싱 실패: %w", err)
	}

	prevPrice, err := decimal.NewFromString(kisPrice.Base)
	if err != nil {
		return nil, fmt.Errorf("전일가 파싱 실패: %w", err)
	}

	// 원환산당일대비 (원화 기준)
	change, err := decimal.NewFromString(kisPrice.TXdif)
	if err != nil {
		return nil, fmt.Errorf("가격변동 파싱 실패: %w", err)
	}

	// 원환산당일등락률 (원화 기준)
	changeRate, err := decimal.NewFromString(kisPrice.TXrat)
	if err != nil {
		return nil, fmt.Errorf("가격변동률 파싱 실패: %w", err)
	}

	// 거래량
	volume := int64(0)
	if kisPrice.Tvol != "" {
		if vol, err := strconv.ParseInt(kisPrice.Tvol, 10, 64); err == nil {
			volume = vol
		}
	}

	// 시가총액
	marketCap := decimal.Zero
	if kisPrice.Tomv != "" {
		if mc, err := decimal.NewFromString(kisPrice.Tomv); err == nil {
			marketCap = mc
		}
	}

	// 고가, 저가, 시가 파싱
	high := decimal.Zero
	if kisPrice.High != "" {
		if h, err := decimal.NewFromString(kisPrice.High); err == nil {
			high = h
		}
	}

	low := decimal.Zero
	if kisPrice.Low != "" {
		if l, err := decimal.NewFromString(kisPrice.Low); err == nil {
			low = l
		}
	}

	open := decimal.Zero
	if kisPrice.Open != "" {
		if o, err := decimal.NewFromString(kisPrice.Open); err == nil {
			open = o
		}
	}

	return &portfolio.StockPrice{
		Symbol:        kisPrice.Rsym,
		Price:         price,
		Change:        change,
		ChangeRate:    changeRate,
		Volume:        volume,
		MarketCap:     marketCap,
		High:          high,
		Low:           low,
		Open:          open,
		PreviousClose: prevPrice,
		Timestamp:     time.Now(),
	}, nil
}

// convertToDataPosition KIS 잔고 데이터를 portfolio.Position으로 변환
func (d *DataAdapter) convertToDataPosition(kisBalance KISBalanceOutput1, userID string) (*portfolio.Position, error) {
	// 문자열을 decimal로 변환
	quantity, err := decimal.NewFromString(kisBalance.OvrsCblcQty)
	if err != nil {
		return nil, fmt.Errorf("잔고수량 파싱 실패: %w", err)
	}

	avgPrice, err := decimal.NewFromString(kisBalance.PchsAvgPric)
	if err != nil {
		return nil, fmt.Errorf("매입평균가격 파싱 실패: %w", err)
	}

	currentPrice, err := decimal.NewFromString(kisBalance.NowPric2)
	if err != nil {
		return nil, fmt.Errorf("현재가격 파싱 실패: %w", err)
	}

	profitAmount, err := decimal.NewFromString(kisBalance.FrcrEvluPflsAmt)
	if err != nil {
		return nil, fmt.Errorf("평가손익금액 파싱 실패: %w", err)
	}

	profitRate, err := decimal.NewFromString(kisBalance.EvluPflsRt)
	if err != nil {
		return nil, fmt.Errorf("평가손익률 파싱 실패: %w", err)
	}

	// 총 가치 계산
	totalValue := quantity.Mul(currentPrice)

	// 일일 수익은 별도 API로 조회 필요
	dailyProfit := decimal.Zero
	dailyProfitRate := decimal.Zero

	return &portfolio.Position{
		UserID:          userID,
		Symbol:          kisBalance.OvrsPdno,
		CompanyName:     kisBalance.OvrsItemName,
		Quantity:        quantity,
		AveragePrice:    avgPrice,
		CurrentPrice:    currentPrice,
		TotalValue:      totalValue,
		TotalProfit:     profitAmount,
		ProfitRate:      profitRate,
		DailyProfit:     dailyProfit,
		DailyProfitRate: dailyProfitRate,
		UpdatedAt:       time.Now(),
	}, nil
}
