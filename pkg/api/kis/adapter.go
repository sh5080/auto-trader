package kis

import (
	"fmt"
	"strconv"
	"time"

	"auto-trader/pkg/domain/portfolio"

	"github.com/shopspring/decimal"
)

// Adapter KIS API 응답을 도메인 모델로 변환하는 어댑터
type Adapter struct{}

// NewAdapter 새로운 어댑터 생성
func NewAdapter() *Adapter {
	return &Adapter{}
}

// ConvertBalanceToPosition KIS 잔고 데이터를 Position으로 변환
func (a *Adapter) ConvertBalanceToPosition(kisBalance KISBalanceOutput1) (*portfolio.Position, error) {
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
		LastUpdated:     time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

// ConvertPriceToStockPrice KIS 현재가 데이터를 StockPrice로 변환
func (a *Adapter) ConvertPriceToStockPrice(kisPrice KISPriceOutput) (*portfolio.StockPrice, error) {
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

// ConvertBalanceToPortfolio KIS 잔고 데이터를 Portfolio로 변환
func (a *Adapter) ConvertBalanceToPortfolio(userID string, balanceResp *KISBalanceResponse) (*portfolio.Portfolio, error) {
	if len(balanceResp.Output1) == 0 {
		return nil, fmt.Errorf("잔고 데이터가 없습니다")
	}

	// 포지션들 변환
	totalValue := decimal.Zero
	totalProfit := decimal.Zero

	for _, kisBalance := range balanceResp.Output1 {
		position, err := a.ConvertBalanceToPosition(kisBalance)
		if err != nil {
			return nil, fmt.Errorf("포지션 변환 실패: %w", err)
		}
		position.UserID = userID

		totalValue = totalValue.Add(position.TotalValue)
		totalProfit = totalProfit.Add(position.TotalProfit)
	}

	// 전체 수익률 계산
	profitRate := decimal.Zero
	if totalValue.GreaterThan(decimal.Zero) {
		profitRate = totalProfit.Div(totalValue).Mul(decimal.NewFromInt(100))
	}

	return &portfolio.Portfolio{
		UserID:      userID,
		TotalValue:  totalValue,
		TotalProfit: totalProfit,
		ProfitRate:  profitRate,
		LastUpdated: time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
