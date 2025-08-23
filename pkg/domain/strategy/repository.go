package strategy

import (
	"context"
	"fmt"

	"auto-trader/ent"
	"auto-trader/ent/strategy"
	"auto-trader/pkg/domain/strategy/dto"

	"github.com/google/uuid"
)

// Repository 전략 데이터 접근 인터페이스
type Repository interface {
	// 기본 CRUD
	Create(input dto.CreateStrategyBody) (*ent.Strategy, error)
	GetByID(id uuid.UUID) (*ent.Strategy, error)
	Update(id uuid.UUID, input dto.UpdateStrategyBody) (*ent.Strategy, error)
	Delete(id uuid.UUID) error

	// 기본 조회
	GetAll(limit, offset int) ([]*ent.Strategy, error)
	Count() (int, error)

	// 전략 특화 메서드
	GetByUserID(userID uuid.UUID, limit, offset int) ([]*ent.Strategy, error)
	GetBySymbol(symbol string) ([]*ent.Strategy, error)
	GetActiveStrategies() ([]*ent.Strategy, error)

	// 관계 조회
	GetStrategyWithUser(id uuid.UUID) (*ent.Strategy, error)
	GetStrategyWithExecutions(id uuid.UUID) (*ent.Strategy, error)

	// 통계
	CountByUser(userID uuid.UUID) (int, error)
	CountBySymbol(symbol string) (int, error)
}

// EntRepository ent 기반 구현체
type EntRepository struct {
	client *ent.Client
}

// NewEntRepository ent 기반 Repository 생성
func NewEntRepository(client *ent.Client) Repository {
	return &EntRepository{client: client}
}

// 헬퍼 함수들
func (r *EntRepository) getContext() context.Context {
	return context.Background()
}

// Create 전략 생성
func (r *EntRepository) Create(input dto.CreateStrategyBody) (*ent.Strategy, error) {
	strategy, err := r.client.Strategy.Create().
		SetName(input.Name).
		SetSymbol(input.Symbol).
		SetDescription(*input.Description).
		SetUserID(input.UserID).
		SetActive(input.Active).
		Save(r.getContext())

	if err != nil {
		return nil, fmt.Errorf("failed to create strategy: %w", err)
	}

	return strategy, nil
}

// GetByID ID로 전략 조회
func (r *EntRepository) GetByID(id uuid.UUID) (*ent.Strategy, error) {
	strategy, err := r.client.Strategy.Get(r.getContext(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get strategy by id: %w", err)
	}
	return strategy, nil
}

// Update 전략 정보 수정
func (r *EntRepository) Update(id uuid.UUID, input dto.UpdateStrategyBody) (*ent.Strategy, error) {
	updateQuery := r.client.Strategy.UpdateOneID(id)

	if input.Name != nil {
		updateQuery.SetName(*input.Name)
	}
	if input.Symbol != nil {
		updateQuery.SetSymbol(*input.Symbol)
	}
	if input.Description != nil {
		updateQuery.SetDescription(*input.Description)
	}
	if input.Active != nil {
		updateQuery.SetActive(*input.Active)
	}

	strategy, err := updateQuery.Save(r.getContext())
	if err != nil {
		return nil, fmt.Errorf("failed to update strategy: %w", err)
	}

	return strategy, nil
}

// Delete 전략 삭제
func (r *EntRepository) Delete(id uuid.UUID) error {
	err := r.client.Strategy.DeleteOneID(id).Exec(r.getContext())
	if err != nil {
		return fmt.Errorf("failed to delete strategy: %w", err)
	}
	return nil
}

// GetAll 모든 전략 조회 (페이지네이션)
func (r *EntRepository) GetAll(limit, offset int) ([]*ent.Strategy, error) {
	strategies, err := r.client.Strategy.Query().
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(strategy.FieldCreatedAt)).
		All(r.getContext())

	if err != nil {
		return nil, fmt.Errorf("failed to get strategies: %w", err)
	}
	return strategies, nil
}

// Count 전체 전략 수
func (r *EntRepository) Count() (int, error) {
	count, err := r.client.Strategy.Query().Count(r.getContext())
	if err != nil {
		return 0, fmt.Errorf("failed to count strategies: %w", err)
	}
	return count, nil
}

// GetByUserID 사용자별 전략 조회
func (r *EntRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]*ent.Strategy, error) {
	strategies, err := r.client.Strategy.Query().
		Where(strategy.UserID(userID)).
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(strategy.FieldCreatedAt)).
		All(r.getContext())

	if err != nil {
		return nil, fmt.Errorf("failed to get strategies by user: %w", err)
	}
	return strategies, nil
}

// GetBySymbol 심볼별 전략 조회
func (r *EntRepository) GetBySymbol(symbol string) ([]*ent.Strategy, error) {
	strategies, err := r.client.Strategy.Query().
		Where(strategy.Symbol(symbol)).
		Order(ent.Desc(strategy.FieldCreatedAt)).
		All(r.getContext())

	if err != nil {
		return nil, fmt.Errorf("failed to get strategies by symbol: %w", err)
	}
	return strategies, nil
}

// GetActiveStrategies 활성 전략만 조회
func (r *EntRepository) GetActiveStrategies() ([]*ent.Strategy, error) {
	strategies, err := r.client.Strategy.Query().
		Where(strategy.Active(true)).
		Order(ent.Desc(strategy.FieldCreatedAt)).
		All(r.getContext())

	if err != nil {
		return nil, fmt.Errorf("failed to get active strategies: %w", err)
	}
	return strategies, nil
}

// GetStrategyWithUser 전략과 사용자 정보 함께 조회
func (r *EntRepository) GetStrategyWithUser(id uuid.UUID) (*ent.Strategy, error) {
	strategy, err := r.client.Strategy.Query().
		Where(strategy.ID(id)).
		WithUser().
		Only(r.getContext())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get strategy with user: %w", err)
	}
	return strategy, nil
}

// GetStrategyWithExecutions 전략과 실행 정보 함께 조회
func (r *EntRepository) GetStrategyWithExecutions(id uuid.UUID) (*ent.Strategy, error) {
	strategy, err := r.client.Strategy.Query().
		Where(strategy.ID(id)).
		WithExecutions().
		Only(r.getContext())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get strategy with executions: %w", err)
	}
	return strategy, nil
}

// CountByUser 사용자별 전략 수
func (r *EntRepository) CountByUser(userID uuid.UUID) (int, error) {
	count, err := r.client.Strategy.Query().
		Where(strategy.UserID(userID)).
		Count(r.getContext())

	if err != nil {
		return 0, fmt.Errorf("failed to count strategies by user: %w", err)
	}
	return count, nil
}

// CountBySymbol 심볼별 전략 수
func (r *EntRepository) CountBySymbol(symbol string) (int, error) {
	count, err := r.client.Strategy.Query().
		Where(strategy.Symbol(symbol)).
		Count(r.getContext())

	if err != nil {
		return 0, fmt.Errorf("failed to count strategies by symbol: %w", err)
	}
	return count, nil
}
