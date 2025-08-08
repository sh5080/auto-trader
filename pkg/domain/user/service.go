package user

import (
	"auto-trader/pkg/domain/user/dto"
	"auto-trader/pkg/shared/utils"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(dto dto.CreateUserRequest) (*UserExceptPassword, error)
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	VerifyPassword(hashed, password string) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service { return &ServiceImpl{repo: repo} }

func (s *ServiceImpl) CreateUser(dto dto.CreateUserRequest) (*UserExceptPassword, error) {
	log.Println("dto", dto)
	existingUser, err := s.GetByEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	log.Println("existingUser", existingUser)
	if existingUser != nil {
		return nil, utils.Conflict("email", "이미 존재하는 이메일입니다")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := UserExceptPassword{
		ID:        uuid.NewString(),
		Name:      dto.Name,
		Nickname:  dto.Nickname,
		Email:     dto.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.repo.Create(&u, string(hash)); err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *ServiceImpl) GetByID(id string) (*User, error) {
	u, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *ServiceImpl) GetByEmail(email string) (*User, error) {
	u, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *ServiceImpl) VerifyPassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
