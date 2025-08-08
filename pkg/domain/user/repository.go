package user

import (
	"auto-trader/pkg/shared/query"
	"auto-trader/pkg/shared/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(u *UserExceptPassword, password string) error
	GetByID(id string) (*User, error)
	GetByIDExceptPassword(id string) (*UserExceptPassword, error)
	GetByEmail(email string) (*User, error)
}

// DBRepository SQL 기반 구현체
type DBRepository struct {
	db *sqlx.DB
}

func NewDBRepository(db *sqlx.DB) Repository {
	return &DBRepository{db: db}
}

func (r *DBRepository) Create(u *UserExceptPassword, password string) error {
	if u == nil {
		return errors.New("user is nil")
	}
	query := `
        INSERT INTO users (` + query.UsersColumns + `)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := r.db.Exec(query, u.ID, u.Name, u.Nickname, u.Email, password, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return utils.Internal("failed to create user", err)
	}
	return nil
}

func (r *DBRepository) GetByID(id string) (*User, error) {
	var u User
	query := `
        SELECT ` + query.UsersColumns + `
        FROM users WHERE id = $1
    `
	if err := r.db.Get(&u, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &u, nil
}

func (r *DBRepository) GetByIDExceptPassword(id string) (*UserExceptPassword, error) {
	var u UserExceptPassword
	query := `
		SELECT ` + query.UsersExceptPasswordColumns + `
		FROM users WHERE id = $1
	`
	if err := r.db.Get(&u, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &u, nil
}
func (r *DBRepository) GetByEmail(email string) (*User, error) {
	var u User
	// 이메일 대소문자 무시
	query := `
        SELECT ` + query.UsersColumns + `
        FROM users WHERE LOWER(email) = $1
    `
	if err := r.db.Get(&u, query, strings.ToLower(email)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &u, nil
}
