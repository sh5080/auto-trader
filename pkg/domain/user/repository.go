package user

import (
	"context"
	"fmt"

	"auto-trader/ent"
	"auto-trader/ent/user"
	"auto-trader/pkg/domain/user/dto"

	"github.com/google/uuid"
)

// Repository 사용자 데이터 접근 인터페이스
type Repository interface {
	// 기본 CRUD
	Create(input dto.CreateUserBody) (*ent.User, error)
	GetByID(id uuid.UUID, includePassword bool) (*ent.User, error)
	Update(id uuid.UUID, input dto.UpdateUserBody) (*ent.User, error)
	Delete(id uuid.UUID) error

	// 기본 조회
	GetAll(limit, offset int, includePassword bool) ([]*ent.User, error)
	Count() (int, error)

	// 사용자 특화 메서드
	GetByEmail(email string, includePassword bool) (*ent.User, error)
	GetByNickname(nickname string, includePassword bool) (*ent.User, error)
	GetActiveUsers(includePassword bool) ([]*ent.User, error)

	// 관계 조회
	GetUserWithStrategies(id uuid.UUID, includePassword bool) (*ent.User, error)
	GetUserWithPortfolios(id uuid.UUID, includePassword bool) (*ent.User, error)

	// 통계
	CountByStatus(isValid bool) (int, error)
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

// 패스워드 제외 필드 선택
func (r *EntRepository) selectFieldsWithoutPassword() []string {
	return []string{
		user.FieldID,
		user.FieldName,
		user.FieldNickname,
		user.FieldEmail,
		user.FieldIsValid,
		user.FieldCreatedAt,
		user.FieldUpdatedAt,
	}
}

// 사용자 조회 헬퍼 (패스워드 포함 여부에 따라)
func (r *EntRepository) getUserWithPasswordOption(query *ent.UserQuery, includePassword bool) (*ent.User, error) {
	ctx := r.getContext()

	if includePassword {
		return query.Only(ctx)
	} else {
		return query.Select(r.selectFieldsWithoutPassword()...).Only(ctx)
	}
}

// 사용자 목록 조회 헬퍼 (패스워드 포함 여부에 따라)
func (r *EntRepository) getUsersWithPasswordOption(query *ent.UserQuery, includePassword bool) ([]*ent.User, error) {
	ctx := r.getContext()

	if includePassword {
		return query.All(ctx)
	} else {
		return query.Select(r.selectFieldsWithoutPassword()...).All(ctx)
	}
}

// Create 사용자 생성
func (r *EntRepository) Create(input dto.CreateUserBody) (*ent.User, error) {
	user, err := r.client.User.Create().
		SetName(input.Name).
		SetNickname(input.Nickname).
		SetEmail(input.Email).
		SetPassword(input.Password).
		SetIsValid(true).
		Save(r.getContext())

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID ID로 사용자 조회
func (r *EntRepository) GetByID(id uuid.UUID, includePassword bool) (*ent.User, error) {
	query := r.client.User.Query().Where(user.ID(id))
	return r.getUserWithPasswordOption(query, includePassword)
}

// GetByEmail 이메일로 사용자 조회
func (r *EntRepository) GetByEmail(email string, includePassword bool) (*ent.User, error) {
	query := r.client.User.Query().Where(user.Email(email))
	return r.getUserWithPasswordOption(query, includePassword)
}

// GetByNickname 닉네임으로 사용자 조회
func (r *EntRepository) GetByNickname(nickname string, includePassword bool) (*ent.User, error) {
	query := r.client.User.Query().Where(user.Nickname(nickname))
	return r.getUserWithPasswordOption(query, includePassword)
}

// GetActiveUsers 활성 사용자만 조회
func (r *EntRepository) GetActiveUsers(includePassword bool) ([]*ent.User, error) {
	query := r.client.User.Query().
		Where(user.IsValid(true)).
		Order(ent.Desc(user.FieldCreatedAt))

	return r.getUsersWithPasswordOption(query, includePassword)
}

// GetUserWithStrategies 사용자와 전략 정보 함께 조회
func (r *EntRepository) GetUserWithStrategies(id uuid.UUID, includePassword bool) (*ent.User, error) {
	query := r.client.User.Query().
		Where(user.ID(id)).
		WithStrategies()

	return r.getUserWithPasswordOption(query, includePassword)
}

// GetUserWithPortfolios 사용자와 포트폴리오 정보 함께 조회
func (r *EntRepository) GetUserWithPortfolios(id uuid.UUID, includePassword bool) (*ent.User, error) {
	query := r.client.User.Query().
		Where(user.ID(id)).
		WithPortfolios()

	return r.getUserWithPasswordOption(query, includePassword)
}

// Update 사용자 정보 수정
func (r *EntRepository) Update(id uuid.UUID, input dto.UpdateUserBody) (*ent.User, error) {
	updateQuery := r.client.User.UpdateOneID(id)

	if input.Name != nil {
		updateQuery.SetName(*input.Name)
	}
	if input.Nickname != nil {
		updateQuery.SetNickname(*input.Nickname)
	}
	if input.Email != nil {
		updateQuery.SetEmail(*input.Email)
	}
	if input.Password != nil {
		updateQuery.SetPassword(*input.Password)
	}
	if input.IsValid != nil {
		updateQuery.SetIsValid(*input.IsValid)
	}

	user, err := updateQuery.Save(r.getContext())
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// Delete 사용자 삭제
func (r *EntRepository) Delete(id uuid.UUID) error {
	err := r.client.User.DeleteOneID(id).Exec(r.getContext())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetAll 모든 사용자 조회 (페이지네이션)
func (r *EntRepository) GetAll(limit, offset int, includePassword bool) ([]*ent.User, error) {
	query := r.client.User.Query().
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(user.FieldCreatedAt))

	return r.getUsersWithPasswordOption(query, includePassword)
}

// Count 전체 사용자 수
func (r *EntRepository) Count() (int, error) {
	count, err := r.client.User.Query().Count(r.getContext())
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}

// CountByStatus 상태별 사용자 수
func (r *EntRepository) CountByStatus(isValid bool) (int, error) {
	count, err := r.client.User.Query().
		Where(user.IsValid(isValid)).
		Count(r.getContext())

	if err != nil {
		return 0, fmt.Errorf("failed to count users by status: %w", err)
	}
	return count, nil
}
