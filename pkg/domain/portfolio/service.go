package portfolio

import (
	"auto-trader/pkg/domain/portfolio/dto"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// Service 포트폴리오 관리 서비스 인터페이스
type Service interface {
	// 포트폴리오 관련
	GetPortfolio(userID string, q dto.GetPortfolioQuery) (*Portfolio, error)
	GetPortfolioSummary(userID string, q dto.GetPortfolioSummaryQuery) (*PortfolioSummary, error)

	// 보유 주식 관련
	GetPositions(userID string, q dto.GetPositionsQuery) ([]*Position, error)
	GetPosition(userID string, q dto.GetPositionsQuery) (*Position, error)

	// 주식 가격 관련
	GetCurrentPrice(q dto.GetCurrentPricesQuery) (*StockPrice, error)
	GetCurrentPrices(q dto.GetCurrentPricesQuery) ([]*StockPrice, error)

	// 거래 내역 관련
	GetTradeHistory(userID string, q dto.GetTradeHistoryQuery) ([]*TradeHistory, error)

	// 회사 정보 관련
	GetCompanyInfo(q dto.SymbolPath) (*CompanyInfo, error)

	// 차트 데이터 관련
	GetChartData(q dto.SymbolPath) ([]*ChartData, error)

	// 실시간 데이터 (향후 WebSocket 구현)
	SubscribeToPriceUpdates(q dto.GetCurrentPricesQuery) (<-chan StockPrice, error)

	// 캐시 관리
	RefreshPortfolio(userID string) error
	RefreshPositions(userID string) error
	RefreshPrices(q dto.GetCurrentPricesQuery) error
}

// ServiceImpl 포트폴리오 서비스 구현체
type ServiceImpl struct {
	repository         Repository
	externalDataSource ExternalDataSource
	cacheConfig        CacheConfig
}

// NewService 새로운 포트폴리오 서비스 생성
func NewService(repository Repository, externalDataSource ExternalDataSource, cacheConfig CacheConfig) Service {
	return &ServiceImpl{
		repository:         repository,
		externalDataSource: externalDataSource,
		cacheConfig:        cacheConfig,
	}
}

// isCacheValid 캐시 유효성 검사
func (s *ServiceImpl) isCacheValid(lastSyncAt time.Time, cacheTTL time.Duration) bool {
	return time.Since(lastSyncAt) < cacheTTL
}

// GetPortfolio 포트폴리오 조회
func (s *ServiceImpl) GetPortfolio(userID string, q dto.GetPortfolioQuery) (*Portfolio, error) {
	// 캐시된 포트폴리오 조회
	if !q.ForceRefresh {
		portfolio, err := s.repository.GetPortfolio(userID)
		if err == nil && s.isCacheValid(portfolio.LastUpdated, s.cacheConfig.PortfolioCacheTTL) {
			logrus.Debugf("캐시된 포트폴리오 반환: %s", userID)
			return portfolio, nil
		}
	}

	// 외부 데이터 소스에서 최신 데이터 조회
	// TODO: 계좌번호를 어떻게 가져올지 결정 필요
	accountNo := "계좌번호" // 실제로는 사용자별 계좌번호 매핑 필요
	positions, err := s.externalDataSource.GetBalance(accountNo)
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

	portfolio := &Portfolio{
		ID:          uuid.New().String(),
		UserID:      userID,
		TotalValue:  totalValue,
		TotalProfit: totalProfit,
		ProfitRate:  profitRate,
		LastUpdated: time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 캐시에 저장
	if err := s.repository.SavePortfolio(portfolio); err != nil {
		logrus.Warnf("포트폴리오 캐시 저장 실패: %v", err)
	}

	return portfolio, nil
}

// GetPortfolioSummary 포트폴리오 요약 조회
func (s *ServiceImpl) GetPortfolioSummary(userID string, q dto.GetPortfolioSummaryQuery) (*PortfolioSummary, error) {
	// 외부 데이터 소스에서 요약 조회
	accountNo := "계좌번호" // 실제로는 사용자별 계좌번호 매핑 필요
	summary, err := s.externalDataSource.GetPortfolioSummary(userID, accountNo)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

// GetPositions 보유 주식 목록 조회
func (s *ServiceImpl) GetPositions(userID string, q dto.GetPositionsQuery) ([]*Position, error) {
	// 캐시된 포지션 조회
	if !q.ForceRefresh {
		positions, err := s.repository.GetPositions(userID)
		if err == nil && len(positions) > 0 {
			// 캐시 유효성 검사 (첫 번째 포지션 기준)
			if s.isCacheValid(positions[0].LastUpdated, s.cacheConfig.PositionCacheTTL) {
				logrus.Debugf("캐시된 포지션 반환: %s", userID)
				return positions, nil
			}
		}
	}

	// 외부 데이터 소스에서 최신 데이터 조회
	accountNo := "계좌번호" // 실제로는 사용자별 계좌번호 매핑 필요
	positions, err := s.externalDataSource.GetBalance(accountNo)
	if err != nil {
		return nil, err
	}

	// 특정 심볼 필터링
	if q.Symbol != "" {
		var filtered []*Position
		for _, pos := range positions {
			if pos.Symbol == q.Symbol {
				filtered = append(filtered, pos)
			}
		}
		positions = filtered
	}

	// 캐시에 저장
	for _, pos := range positions {
		pos.UserID = userID
		if err := s.repository.SavePosition(pos); err != nil {
			logrus.Warnf("포지션 캐시 저장 실패: %v", err)
		}
	}

	return positions, nil
}

// GetPosition 특정 보유 주식 조회
func (s *ServiceImpl) GetPosition(userID string, q dto.GetPositionsQuery) (*Position, error) {
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
func (s *ServiceImpl) GetCurrentPrice(q dto.GetCurrentPricesQuery) (*StockPrice, error) {
	// 캐시된 가격 조회
	if !q.ForceRefresh {
		price, err := s.repository.GetStockPrice(q.Symbols)
		if err == nil && s.isCacheValid(price.Timestamp, s.cacheConfig.PriceCacheTTL) {
			logrus.Debugf("캐시된 가격 반환: %s", q.Symbols)
			return price, nil
		}
	}

	// 외부 데이터 소스에서 최신 데이터 조회
	price, err := s.externalDataSource.GetCurrentPrice(q.Symbols)
	if err != nil {
		return nil, err
	}

	// 캐시에 저장
	if err := s.repository.SaveStockPrice(price); err != nil {
		logrus.Warnf("가격 캐시 저장 실패: %v", err)
	}

	return price, nil
}

// GetCurrentPrices 여러 종목 현재가 조회
func (s *ServiceImpl) GetCurrentPrices(q dto.GetCurrentPricesQuery) ([]*StockPrice, error) {
	return s.externalDataSource.GetCurrentPrices(strings.Split(q.Symbols, ","))
}

// GetDailyProfit 일일 수익 조회
func (s *ServiceImpl) GetDailyProfit(q dto.SymbolPath) (decimal.Decimal, error) {
	// TODO: 일일 수익 계산 로직 구현
	return decimal.Zero, nil
}

// GetTradeHistory 거래 내역 조회
func (s *ServiceImpl) GetTradeHistory(userID string, q dto.GetTradeHistoryQuery) ([]*TradeHistory, error) {
	return s.repository.GetTradeHistory(userID, q)
}

// GetCompanyInfo 회사 정보 조회
func (s *ServiceImpl) GetCompanyInfo(q dto.SymbolPath) (*CompanyInfo, error) {
	// TODO: 외부 API에서 회사 정보 조회 구현
	return nil, fmt.Errorf("not implemented")
}

// GetChartData 차트 데이터 조회
func (s *ServiceImpl) GetChartData(q dto.SymbolPath) ([]*ChartData, error) {
	// TODO: 외부 API에서 차트 데이터 조회 구현
	return nil, fmt.Errorf("not implemented")
}

// SubscribeToPriceUpdates 실시간 가격 업데이트 구독
func (s *ServiceImpl) SubscribeToPriceUpdates(q dto.GetCurrentPricesQuery) (<-chan StockPrice, error) {
	// TODO: WebSocket을 통한 실시간 데이터 구독 구현
	ch := make(chan StockPrice)
	close(ch)
	return ch, nil
}

// RefreshPortfolio 포트폴리오 새로고침
func (s *ServiceImpl) RefreshPortfolio(userID string) error {
	_, err := s.GetPortfolio(userID, dto.GetPortfolioQuery{
		ForceRefresh: true,
	})
	return err
}

// RefreshPositions 포지션 새로고침
func (s *ServiceImpl) RefreshPositions(userID string) error {
	_, err := s.GetPositions(userID, dto.GetPositionsQuery{
		Symbol:       "",
		ForceRefresh: true,
	})
	return err
}

// RefreshPrices 가격 새로고침
func (s *ServiceImpl) RefreshPrices(q dto.GetCurrentPricesQuery) error {
	_, err := s.GetCurrentPrices(q)
	return err
}
