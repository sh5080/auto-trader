package portfolio

import (
	"auto-trader/pkg/domain/portfolio/dto"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Service 포트폴리오 관리 서비스 인터페이스
type Service interface {
	// 포트폴리오 관련
	GetPortfolio(userID string, q dto.GetPortfolioQuery) (*dto.Portfolio, error)
	GetPortfolioSummary(userID string, q dto.GetPortfolioSummaryQuery) (*dto.PortfolioSummary, error)

	// 보유 주식 관련
	GetPositions(userID string, q dto.GetPositionsQuery) ([]*dto.Position, error)
	GetPosition(userID string, q dto.GetPositionsQuery) (*dto.Position, error)

	// 주식 가격 관련
	GetCurrentPrice(q dto.GetCurrentPricesQuery) (*dto.StockPrice, error)
	GetCurrentPrices(q dto.GetCurrentPricesQuery) ([]*dto.StockPrice, error)

	// 거래 내역 관련
	GetTradeHistory(userID string, q dto.GetTradeHistoryQuery) ([]*dto.TradeHistory, error)

	// 회사 정보 관련
	GetCompanyInfo(q dto.SymbolPath) (*dto.CompanyInfo, error)

	// 차트 데이터 관련
	GetChartData(q dto.SymbolPath) ([]*dto.ChartData, error)

	// 실시간 데이터 (향후 WebSocket 구현)
	SubscribeToPriceUpdates(q dto.GetCurrentPricesQuery) (<-chan dto.StockPrice, error)

	// 캐시 관리
	RefreshPortfolio(userID string) error
	RefreshPositions(userID string) error
	RefreshPrices(q dto.GetCurrentPricesQuery) error
}

// ServiceImpl 포트폴리오 서비스 구현체
type ServiceImpl struct {
	repository Repository
}

// NewService 새로운 포트폴리오 서비스 생성
func NewService(repository Repository) Service {
	return &ServiceImpl{
		repository: repository,
	}
}

// GetPortfolio 포트폴리오 조회
func (s *ServiceImpl) GetPortfolio(userID string, q dto.GetPortfolioQuery) (*dto.Portfolio, error) {
	// UUID 변환
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("잘못된 사용자 ID 형식: %w", err)
	}

	// Repository에서 포트폴리오 조회
	portfolios, err := s.repository.GetByUserID(userUUID, 100, 0) // 적절한 limit, offset 설정
	if err != nil {
		return nil, fmt.Errorf("포트폴리오 조회 실패: %w", err)
	}

	if len(portfolios) == 0 {
		return &dto.Portfolio{
			ID:          uuid.New().String(),
			UserID:      userID,
			TotalValue:  decimal.Zero,
			TotalProfit: decimal.Zero,
			ProfitRate:  decimal.Zero,
			LastUpdated: time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}

	// 포트폴리오 계산
	totalValue := decimal.Zero
	totalProfit := decimal.Zero
	for _, pos := range portfolios {
		// ent.Portfolio의 Quantity를 float64로 변환하여 계산
		quantity, _ := pos.Quantity.Float64()
		totalValue = totalValue.Add(decimal.NewFromFloat(quantity))
		// 수익률은 임시로 0으로 설정 (실제로는 가격 정보 필요)
		totalProfit = totalProfit.Add(decimal.Zero)
	}

	profitRate := decimal.Zero
	if totalValue.GreaterThan(decimal.Zero) {
		profitRate = totalProfit.Div(totalValue).Mul(decimal.NewFromInt(100))
	}

	portfolio := &dto.Portfolio{
		ID:          uuid.New().String(),
		UserID:      userID,
		TotalValue:  totalValue,
		TotalProfit: totalProfit,
		ProfitRate:  profitRate,
		LastUpdated: time.Now(),
		UpdatedAt:   time.Now(),
	}

	return portfolio, nil
}

// GetPortfolioSummary 포트폴리오 요약 조회
func (s *ServiceImpl) GetPortfolioSummary(userID string, q dto.GetPortfolioSummaryQuery) (*dto.PortfolioSummary, error) {
	// UUID 변환
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("잘못된 사용자 ID 형식: %w", err)
	}

	// Repository에서 포트폴리오 통계 조회
	count, err := s.repository.CountByUser(userUUID)
	if err != nil {
		return nil, fmt.Errorf("포트폴리오 수 조회 실패: %w", err)
	}

	totalValue, err := s.repository.GetTotalValueByUser(userUUID)
	if err != nil {
		return nil, fmt.Errorf("총 포트폴리오 가치 조회 실패: %w", err)
	}

	summary := &dto.PortfolioSummary{
		TotalPositions: count,
		TotalValue:     decimal.NewFromFloat(totalValue),
		LastUpdated:    time.Now(),
	}

	return summary, nil
}

// GetPositions 보유 주식 목록 조회
func (s *ServiceImpl) GetPositions(userID string, q dto.GetPositionsQuery) ([]*dto.Position, error) {
	// UUID 변환
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("잘못된 사용자 ID 형식: %w", err)
	}

	// Repository에서 포트폴리오 조회
	portfolios, err := s.repository.GetByUserID(userUUID, 100, 0) // 적절한 limit, offset 설정
	if err != nil {
		return nil, fmt.Errorf("보유 주식 조회 실패: %w", err)
	}

	// ent.Portfolio를 Position으로 변환
	var positions []*dto.Position
	for _, portfolio := range portfolios {
		quantity, _ := portfolio.Quantity.Float64()
		position := &dto.Position{
			Symbol:      portfolio.Symbol,
			Quantity:    decimal.NewFromFloat(quantity),
			TotalValue:  decimal.NewFromFloat(quantity), // 임시로 수량과 동일하게 설정
			TotalProfit: decimal.Zero,                   // 임시로 0으로 설정
			LastUpdated: time.Now(),
		}
		positions = append(positions, position)
	}

	// 특정 심볼 필터링
	if q.Symbol != "" {
		var filtered []*dto.Position
		for _, pos := range positions {
			if pos.Symbol == q.Symbol {
				filtered = append(filtered, pos)
			}
		}
		positions = filtered
	}

	return positions, nil
}

// GetPosition 특정 보유 주식 조회
func (s *ServiceImpl) GetPosition(userID string, q dto.GetPositionsQuery) (*dto.Position, error) {
	positions, err := s.GetPositions(userID, q)
	if err != nil {
		return nil, err
	}

	if len(positions) == 0 {
		return nil, fmt.Errorf("보유 주식이 없습니다: %s", q.Symbol)
	}

	return positions[0], nil
}

// GetCurrentPrice 현재가 조회
func (s *ServiceImpl) GetCurrentPrice(q dto.GetCurrentPricesQuery) (*dto.StockPrice, error) {
	// TODO: 외부 API에서 현재가 조회 구현
	// 현재는 임시 데이터 반환
	price := &dto.StockPrice{
		Symbol:    q.Symbols,
		Price:     decimal.NewFromFloat(100.0), // 임시 가격
		Timestamp: time.Now(),
	}

	return price, nil
}

// GetCurrentPrices 여러 종목 현재가 조회
func (s *ServiceImpl) GetCurrentPrices(q dto.GetCurrentPricesQuery) ([]*dto.StockPrice, error) {
	symbols := strings.Split(q.Symbols, ",")
	var prices []*dto.StockPrice

	for _, symbol := range symbols {
		price := &dto.StockPrice{
			Symbol:    symbol,
			Price:     decimal.NewFromFloat(100.0), // 임시 가격
			Timestamp: time.Now(),
		}
		prices = append(prices, price)
	}

	return prices, nil
}

// GetDailyProfit 일일 수익 조회
func (s *ServiceImpl) GetDailyProfit(q dto.SymbolPath) (decimal.Decimal, error) {
	// TODO: 일일 수익 계산 로직 구현
	return decimal.Zero, nil
}

// GetTradeHistory 거래 내역 조회
func (s *ServiceImpl) GetTradeHistory(userID string, q dto.GetTradeHistoryQuery) ([]*dto.TradeHistory, error) {
	// TODO: 거래 내역 조회 로직 구현
	// 현재는 빈 배열 반환
	return []*dto.TradeHistory{}, nil
}

// GetCompanyInfo 회사 정보 조회
func (s *ServiceImpl) GetCompanyInfo(q dto.SymbolPath) (*dto.CompanyInfo, error) {
	// TODO: 외부 API에서 회사 정보 조회 구현
	return nil, fmt.Errorf("not implemented")
}

// GetChartData 차트 데이터 조회
func (s *ServiceImpl) GetChartData(q dto.SymbolPath) ([]*dto.ChartData, error) {
	// TODO: 외부 API에서 차트 데이터 조회 구현
	return nil, fmt.Errorf("not implemented")
}

// SubscribeToPriceUpdates 실시간 가격 업데이트 구독
func (s *ServiceImpl) SubscribeToPriceUpdates(q dto.GetCurrentPricesQuery) (<-chan dto.StockPrice, error) {
	// TODO: WebSocket을 통한 실시간 데이터 구독 구현
	ch := make(chan dto.StockPrice)
	close(ch)
	return ch, nil
}

// RefreshPortfolio 포트폴리오 새로고침
func (s *ServiceImpl) RefreshPortfolio(userID string) error {
	_, err := s.GetPortfolio(userID, dto.GetPortfolioQuery{})
	return err
}

// RefreshPositions 포지션 새로고침
func (s *ServiceImpl) RefreshPositions(userID string) error {
	_, err := s.GetPositions(userID, dto.GetPositionsQuery{
		Symbol: "",
	})
	return err
}

// RefreshPrices 가격 새로고침
func (s *ServiceImpl) RefreshPrices(q dto.GetCurrentPricesQuery) error {
	_, err := s.GetCurrentPrices(q)
	return err
}
