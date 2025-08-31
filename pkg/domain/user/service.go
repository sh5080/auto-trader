package user

import (
	"auto-trader/pkg/domain/user/dto"
	"auto-trader/pkg/shared/utils"

	"auto-trader/ent"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreateUser(input dto.CreateUserBody) (*ent.User, error)
	GetByID(id uuid.UUID, includePassword ...bool) (*ent.User, error)
	GetByEmail(email string, includePassword ...bool) (*ent.User, error)
	VerifyPassword(hashed, password string) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service { return &ServiceImpl{repo: repo} }

func (s *ServiceImpl) CreateUser(input dto.CreateUserBody) (*ent.User, error) {
	existingUser, err := s.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, utils.Conflict("email", "이미 존재하는 이메일입니다")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 패스워드 해시화
	hashedInput := dto.CreateUserBody{
		Name:     input.Name,
		Nickname: input.Nickname,
		Email:    input.Email,
		Password: string(hash),
	}

	user, err := s.repo.Create(hashedInput)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *ServiceImpl) GetByID(id uuid.UUID, includePassword ...bool) (*ent.User, error) {
	include := false
	if len(includePassword) > 0 {
		include = includePassword[0]
	}

	u, err := s.repo.GetByID(id, include)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *ServiceImpl) GetByEmail(email string, includePassword ...bool) (*ent.User, error) {
	include := false
	if len(includePassword) > 0 {
		include = includePassword[0]
	}

	u, err := s.repo.GetByEmail(email, include)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *ServiceImpl) VerifyPassword(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
