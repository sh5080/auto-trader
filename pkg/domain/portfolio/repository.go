package portfolio

import (
	"context"
	"fmt"

	"auto-trader/ent"
	"auto-trader/ent/portfolio"

	"auto-trader/pkg/domain/portfolio/dto"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Repository 포트폴리오 데이터 접근 인터페이스
type Repository interface {
	// 기본 CRUD
	Create(input dto.CreatePortfolioBody) (*ent.Portfolio, error)
	GetByID(id uuid.UUID) (*ent.Portfolio, error)
	Update(id uuid.UUID, input dto.UpdatePortfolioBody) (*ent.Portfolio, error)
	Delete(id uuid.UUID) error

	// 기본 조회
	GetAll(limit, offset int) ([]*ent.Portfolio, error)
	Count() (int, error)

	// 포트폴리오 특화 메서드
	GetByUserID(userID uuid.UUID, limit, offset int) ([]*ent.Portfolio, error)
	GetBySymbol(symbol string) ([]*ent.Portfolio, error)
	GetByUserAndSymbol(userID uuid.UUID, symbol string) (*ent.Portfolio, error)

	// 관계 조회
	GetPortfolioWithUser(id uuid.UUID) (*ent.Portfolio, error)

	// 통계
	CountByUser(userID uuid.UUID) (int, error)
	CountBySymbol(symbol string) (int, error)
	GetTotalValueByUser(userID uuid.UUID) (float64, error)
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

// Create 포트폴리오 생성
func (r *EntRepository) Create(input dto.CreatePortfolioBody) (*ent.Portfolio, error) {
	portfolio, err := r.client.Portfolio.Create().
		SetSymbol(input.Symbol).
		SetQuantity(decimal.NewFromFloat(input.Quantity)).
		SetUserID(input.UserID).
		Save(r.getContext())

	if err != nil {
		return nil, fmt.Errorf("failed to create portfolio: %w", err)
	}

	return portfolio, nil
}

// GetByID ID로 포트폴리오 조회
func (r *EntRepository) GetByID(id uuid.UUID) (*ent.Portfolio, error) {
	portfolio, err := r.client.Portfolio.Get(r.getContext(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get portfolio by id: %w", err)
	}
	return portfolio, nil
}

// Update 포트폴리오 정보 수정
func (r *EntRepository) Update(id uuid.UUID, input dto.UpdatePortfolioBody) (*ent.Portfolio, error) {
	updateQuery := r.client.Portfolio.UpdateOneID(id)

	if input.Symbol != nil {
		updateQuery.SetSymbol(*input.Symbol)
	}
	if input.Quantity != nil {
		updateQuery.SetQuantity(decimal.NewFromFloat(*input.Quantity))
	}

	portfolio, err := updateQuery.Save(r.getContext())
	if err != nil {
		return nil, fmt.Errorf("failed to update portfolio: %w", err)
	}

	return portfolio, nil
}

// Delete 포트폴리오 삭제
func (r *EntRepository) Delete(id uuid.UUID) error {
	err := r.client.Portfolio.DeleteOneID(id).Exec(r.getContext())
	if err != nil {
		return fmt.Errorf("failed to delete portfolio: %w", err)
	}
	return nil
}

// GetAll 모든 포트폴리오 조회 (페이지네이션)
func (r *EntRepository) GetAll(limit, offset int) ([]*ent.Portfolio, error) {
	portfolios, err := r.client.Portfolio.Query().
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(portfolio.FieldCreatedAt)).
		All(r.getContext())

	if err != nil {
		return nil, fmt.Errorf("failed to get portfolios: %w", err)
	}
	return portfolios, nil
}

// Count 전체 포트폴리오 수
func (r *EntRepository) Count() (int, error) {
	count, err := r.client.Portfolio.Query().Count(r.getContext())
	if err != nil {
		return 0, fmt.Errorf("failed to count portfolios: %w", err)
	}
	return count, nil
}

// GetByUserID 사용자별 포트폴리오 조회
func (r *EntRepository) GetByUserID(userID uuid.UUID, limit, offset int) ([]*ent.Portfolio, error) {
	portfolios, err := r.client.Portfolio.Query().
		Where(portfolio.UserID(userID)).
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(portfolio.FieldCreatedAt)).
		All(r.getContext())

	if err != nil {
		return nil, fmt.Errorf("failed to get portfolios by user: %w", err)
	}
	return portfolios, nil
}

// GetBySymbol 심볼별 포트폴리오 조회
func (r *EntRepository) GetBySymbol(symbol string) ([]*ent.Portfolio, error) {
	portfolios, err := r.client.Portfolio.Query().
		Where(portfolio.Symbol(symbol)).
		Order(ent.Desc(portfolio.FieldCreatedAt)).
		All(r.getContext())

	if err != nil {
		return nil, fmt.Errorf("failed to get portfolios by symbol: %w", err)
	}
	return portfolios, nil
}

// GetByUserAndSymbol 사용자와 심볼로 포트폴리오 조회
func (r *EntRepository) GetByUserAndSymbol(userID uuid.UUID, symbol string) (*ent.Portfolio, error) {
	portfolio, err := r.client.Portfolio.Query().
		Where(
			portfolio.And(
				portfolio.UserID(userID),
				portfolio.Symbol(symbol),
			),
		).
		Only(r.getContext())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get portfolio by user and symbol: %w", err)
	}
	return portfolio, nil
}

// GetPortfolioWithUser 포트폴리오와 사용자 정보 함께 조회
func (r *EntRepository) GetPortfolioWithUser(id uuid.UUID) (*ent.Portfolio, error) {
	portfolio, err := r.client.Portfolio.Query().
		Where(portfolio.ID(id)).
		WithUser().
		Only(r.getContext())

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get portfolio with user: %w", err)
	}
	return portfolio, nil
}

// CountByUser 사용자별 포트폴리오 수
func (r *EntRepository) CountByUser(userID uuid.UUID) (int, error) {
	count, err := r.client.Portfolio.Query().
		Where(portfolio.UserID(userID)).
		Count(r.getContext())

	if err != nil {
		return 0, fmt.Errorf("failed to count portfolios by user: %w", err)
	}
	return count, nil
}

// CountBySymbol 심볼별 포트폴리오 수
func (r *EntRepository) CountBySymbol(symbol string) (int, error) {
	count, err := r.client.Portfolio.Query().
		Where(portfolio.Symbol(symbol)).
		Count(r.getContext())

	if err != nil {
		return 0, fmt.Errorf("failed to count portfolios by symbol: %w", err)
	}
	return count, nil
}

// GetTotalValueByUser 사용자별 총 포트폴리오 가치
func (r *EntRepository) GetTotalValueByUser(userID uuid.UUID) (float64, error) {
	portfolios, err := r.client.Portfolio.Query().
		Where(portfolio.UserID(userID)).
		All(r.getContext())

	if err != nil {
		return 0, fmt.Errorf("failed to get portfolios for total value calculation: %w", err)
	}

	var totalValue float64
	for _, p := range portfolios {
		// decimal.Decimal을 float64로 변환
		quantity, _ := p.Quantity.Float64()
		totalValue += quantity
	}

	return totalValue, nil
}
