package strategy

import (
	"context"
	"fmt"
	"sync"
	"time"

	"auto-trader/pkg/shared/config"
	"auto-trader/pkg/shared/middleware"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// 임시 인터페이스들 (data, order 도메인이 완성되면 제거)
type Collector interface {
	StartPriceStream(symbols []string)
	Stop()
	GetCurrentPrice(symbol string) (*PriceData, error)
	GetDailyProfit(symbol string) (decimal.Decimal, error)
}

type Executor interface {
	// TODO: order 도메인 완성 시 실제 메서드 정의
	ExecuteOrder(ctx context.Context, symbol, action string, quantity, price decimal.Decimal, orderType string) (string, error)
}

// PriceData 가격 데이터 구조체
type PriceData struct {
	Price decimal.Decimal
}

// Service 전략 서비스 인터페이스
type Service interface {
	// 기존 CRUD 메서드들
	GetAllStrategies() ([]*StrategyDetails, error)
	GetStrategy(id string) (*StrategyDetails, error)
	GetStrategyStatus(id string) (*StrategyStatus, error)
	StartStrategy(id string) error
	StopStrategy(id string) error
	RestartStrategy(id string) error
	CreateStrategy(req *CreateStrategyRequest) (*StrategyDetails, error)
	UpdateStrategy(id string, req *UpdateStrategyRequest) (*StrategyDetails, error)
	DeleteStrategy(id string) error
	GetStrategyPerformance(id string) (*StrategyPerformance, error)

	// Manager에서 이전한 메서드들
	Start() error
	Stop() error
	RegisterStrategy(strategy Strategy) error
}

// ServiceImpl 전략 서비스 구현체
type ServiceImpl struct {
	// 리포지토리 (DB 접근)
	repository Repository

	// 외부 의존성들
	dataCollector Collector
	executor      Executor
	riskManager   *middleware.Manager
	config        *config.Config

	// 메모리 상태 관리 (런타임 전략 인스턴스)
	strategies       map[string]Strategy
	activeStrategies map[string]bool
	mutex            sync.RWMutex
	stopChan         chan struct{}
	isRunning        bool
}

// NewService 새로운 전략 서비스 생성
func NewService(
	repository Repository,
	dataCollector Collector,
	executor Executor,
	riskManager *middleware.Manager,
	config *config.Config,
) Service {
	return &ServiceImpl{
		repository:       repository,
		dataCollector:    dataCollector,
		executor:         executor,
		riskManager:      riskManager,
		config:           config,
		strategies:       make(map[string]Strategy),
		activeStrategies: make(map[string]bool),
		stopChan:         make(chan struct{}),
		isRunning:        false,
	}
}

// Start 전략 서비스 시작
func (s *ServiceImpl) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isRunning {
		return nil
	}

	// 기본 전략들 등록
	s.registerDefaultStrategies()

	// 가격 스트림 시작
	if s.dataCollector != nil {
		symbols := s.getAllSymbols()
		if len(symbols) > 0 {
			go s.dataCollector.StartPriceStream(symbols)
		}
	}

	// 전략 실행 루프 시작
	s.stopChan = make(chan struct{})
	go s.strategyLoop()

	s.isRunning = true
	logrus.Info("🚀 전략 서비스 시작됨")
	return nil
}

// Stop 전략 서비스 중지
func (s *ServiceImpl) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isRunning {
		return nil
	}

	close(s.stopChan)
	if s.dataCollector != nil {
		s.dataCollector.Stop()
	}

	s.isRunning = false
	logrus.Info("⏹️  전략 서비스 중지됨")
	return nil
}

// RegisterStrategy 전략 등록
func (s *ServiceImpl) RegisterStrategy(strategy Strategy) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.strategies[strategy.ID()] = strategy
	logrus.Infof("📝 전략 등록: %s (%s)", strategy.Name(), strategy.ID())
	return nil
}

func (s *ServiceImpl) strategyLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.executeActiveStrategies()
		case <-s.stopChan:
			return
		}
	}
}

func (s *ServiceImpl) executeActiveStrategies() {
	s.mutex.RLock()
	activeStrategies := make(map[string]Strategy)
	for id, strategy := range s.strategies {
		if s.activeStrategies[id] {
			activeStrategies[id] = strategy
		}
	}
	s.mutex.RUnlock()

	for id, strategy := range activeStrategies {
		go func(strategyID string, strat Strategy) {
			if err := strat.Execute(); err != nil {
				logrus.Errorf("❌ 전략 실행 오류 (%s): %v", strategyID, err)
			}
		}(id, strategy)
	}
}

func (s *ServiceImpl) getAllSymbols() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	symbolMap := make(map[string]bool)
	for _, strategy := range s.strategies {
		for _, symbol := range strategy.Symbols() {
			symbolMap[symbol] = true
		}
	}

	var symbols []string
	for symbol := range symbolMap {
		symbols = append(symbols, symbol)
	}
	return symbols
}

// registerDefaultStrategies DB에서 활성 전략들을 동적으로 로드
func (s *ServiceImpl) registerDefaultStrategies() {
	// 현재는 nil 체크로 안전하게 처리 (향후 의존성 완성 시 활성화)
	if s.dataCollector == nil || s.executor == nil {
		logrus.Warn("⚠️  전략 의존성이 완전하지 않음 - 기본 전략 등록 스킵")
		return
	}

	// DB에서 활성 전략들 조회
	strategies, err := s.repository.GetAll()
	if err != nil {
		logrus.Errorf("❌ 활성 전략 조회 실패: %v", err)
		return
	}

	// 활성 전략들을 동적 전략으로 등록
	activeCount := 0
	for _, strategy := range strategies {
		if strategy.Active {
			// 전략 설정 로드
			config, err := s.repository.GetConfig(strategy.ID)
			if err != nil {
				logrus.Errorf("❌ 전략 설정 로드 실패 (%s): %v", strategy.ID, err)
				continue
			}

			// 동적 전략 생성
			dynamicStrategy := NewDynamicStrategy(
				s.dataCollector,
				s.executor,
				s.riskManager,
				s.config,
				config,
			)

			// 전략 등록
			if err := s.RegisterStrategy(dynamicStrategy); err != nil {
				logrus.Errorf("❌ 동적 전략 등록 실패 (%s): %v", strategy.ID, err)
				continue
			}

			logrus.Infof("✅ 동적 전략 등록: %s (%s)", strategy.Name, strategy.ID)
			activeCount++
		}
	}

	logrus.Infof("🎯 총 %d개의 활성 전략이 동적으로 로드되었습니다", activeCount)
}

// GetAllStrategies 모든 전략 조회 (Repository 활용)
func (s *ServiceImpl) GetAllStrategies() ([]*StrategyDetails, error) {
	strategies, err := s.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("전략 목록 조회 실패: %w", err)
	}

	// 각 전략의 성과 정보 추가
	for _, strategy := range strategies {
		if performance, err := s.repository.GetPerformance(strategy.ID); err == nil {
			strategy.Performance = performance
		}
	}

	return strategies, nil
}

// GetStrategy 특정 전략 조회 (Repository 활용)
func (s *ServiceImpl) GetStrategy(id string) (*StrategyDetails, error) {
	strategy, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("전략 조회 실패: %w", err)
	}

	// 성과 정보 추가
	if performance, err := s.repository.GetPerformance(id); err == nil {
		strategy.Performance = performance
	}

	return strategy, nil
}

// GetStrategyStatus 전략 상태 조회 (Repository 활용)
func (s *ServiceImpl) GetStrategyStatus(id string) (*StrategyStatus, error) {
	// 전략 존재 확인
	_, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("전략을 찾을 수 없습니다: %w", err)
	}

	// 상태 조회
	status, err := s.repository.GetStatus(id)
	if err != nil {
		return nil, fmt.Errorf("전략 상태 조회 실패: %w", err)
	}

	return status, nil
}

// StartStrategy 전략 시작 (Repository 활용)
func (s *ServiceImpl) StartStrategy(id string) error {
	// 전략 존재 확인
	strategy, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("전략을 찾을 수 없음: %w", err)
	}

	// 상태 업데이트
	status := &StrategyStatus{
		ID:             id,
		Status:         "running",
		LastExecution:  time.Now(),
		ExecutionCount: 0,
		Uptime:         0,
	}

	if err := s.repository.SaveStatus(status); err != nil {
		return fmt.Errorf("전략 상태 저장 실패: %w", err)
	}

	// 전략 활성화 (DB 업데이트)
	strategy.Active = true
	if err := s.repository.Update(strategy); err != nil {
		return fmt.Errorf("전략 활성화 실패: %w", err)
	}

	// 메모리 상태 업데이트
	s.mutex.Lock()
	if strategyInstance, exists := s.strategies[id]; exists {
		s.activeStrategies[id] = true
		strategyInstance.Start()
	}
	s.mutex.Unlock()

	logrus.Infof("▶️  전략 시작: %s (%s)", strategy.Name, id)
	return nil
}

// StopStrategy 전략 중지 (Repository 활용)
func (s *ServiceImpl) StopStrategy(id string) error {
	// 전략 존재 확인
	strategy, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("전략을 찾을 수 없음: %w", err)
	}

	// 상태 업데이트
	status := &StrategyStatus{
		ID:             id,
		Status:         "stopped",
		LastExecution:  time.Now(),
		ExecutionCount: 0,
		Uptime:         0,
	}

	if err := s.repository.SaveStatus(status); err != nil {
		return fmt.Errorf("전략 상태 저장 실패: %w", err)
	}

	// 전략 비활성화 (DB 업데이트)
	strategy.Active = false
	if err := s.repository.Update(strategy); err != nil {
		return fmt.Errorf("전략 비활성화 실패: %w", err)
	}

	// 메모리 상태 업데이트
	s.mutex.Lock()
	if strategyInstance, exists := s.strategies[id]; exists {
		delete(s.activeStrategies, id)
		strategyInstance.Stop()
	}
	s.mutex.Unlock()

	logrus.Infof("⏸️  전략 중지: %s (%s)", strategy.Name, id)
	return nil
}

// RestartStrategy 전략 재시작
func (s *ServiceImpl) RestartStrategy(id string) error {
	if err := s.StopStrategy(id); err != nil {
		return err
	}

	// 잠시 대기
	time.Sleep(1 * time.Second)

	return s.StartStrategy(id)
}

// GetStrategyPerformance 전략 성과 조회 (Repository 활용)
func (s *ServiceImpl) GetStrategyPerformance(id string) (*StrategyPerformance, error) {
	// 전략 존재 확인
	_, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("전략을 찾을 수 없습니다: %w", err)
	}

	// 성과 조회
	performance, err := s.repository.GetPerformance(id)
	if err != nil {
		return nil, fmt.Errorf("전략 성과 조회 실패: %w", err)
	}

	return performance, nil
}

// CreateStrategy 새로운 전략 생성 (Repository 활용)
func (s *ServiceImpl) CreateStrategy(req *CreateStrategyRequest) (*StrategyDetails, error) {
	// 전략 생성
	strategy := &StrategyDetails{
		Name:        req.Name,
		Description: req.Description,
		Symbols:     req.Symbols,
		Active:      false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// DB에 저장
	if err := s.repository.Create(strategy); err != nil {
		return nil, fmt.Errorf("전략 생성 실패: %w", err)
	}

	logrus.Infof("✅ 전략 생성 완료: %s (%s)", strategy.Name, strategy.ID)
	return strategy, nil
}

// UpdateStrategy 전략 수정 (Repository 활용)
func (s *ServiceImpl) UpdateStrategy(id string, req *UpdateStrategyRequest) (*StrategyDetails, error) {
	// 기존 전략 조회
	strategy, err := s.repository.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("전략을 찾을 수 없습니다: %w", err)
	}

	// 업데이트할 필드 적용
	if req.Name != "" {
		strategy.Name = req.Name
	}
	if req.Description != "" {
		strategy.Description = req.Description
	}
	if len(req.Symbols) > 0 {
		strategy.Symbols = req.Symbols
	}
	strategy.UpdatedAt = time.Now()

	// DB에 저장
	if err := s.repository.Update(strategy); err != nil {
		return nil, fmt.Errorf("전략 수정 실패: %w", err)
	}

	logrus.Infof("✅ 전략 수정 완료: %s (%s)", strategy.Name, strategy.ID)
	return strategy, nil
}

// DeleteStrategy 전략 삭제 (Repository 활용)
func (s *ServiceImpl) DeleteStrategy(id string) error {
	// 전략 존재 확인
	strategy, err := s.repository.GetByID(id)
	if err != nil {
		return fmt.Errorf("전략을 찾을 수 없습니다: %w", err)
	}

	// 실행 중인 전략이면 먼저 중지
	if strategy.Active {
		if err := s.StopStrategy(id); err != nil {
			logrus.Warnf("⚠️  전략 중지 실패, 강제 삭제 진행: %v", err)
		}
	}

	// DB에서 삭제
	if err := s.repository.Delete(id); err != nil {
		return fmt.Errorf("전략 삭제 실패: %w", err)
	}

	// 메모리에서도 제거
	s.mutex.Lock()
	delete(s.strategies, id)
	delete(s.activeStrategies, id)
	s.mutex.Unlock()

	logrus.Infof("🗑️  전략 삭제 완료: %s (%s)", strategy.Name, id)
	return nil
}
