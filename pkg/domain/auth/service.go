package auth

import (
	"errors"
	"strings"
	"time"

	"auto-trader/pkg/domain/auth/dto"
	"auto-trader/pkg/domain/user"
	"auto-trader/pkg/shared/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service interface {
	Login(dto dto.LoginBody) (*TokenPair, error)
	Refresh(refreshAuthHeader string) (*TokenPair, error)
}

type ServiceImpl struct {
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
	users      user.Service
}

func NewService(secret string, accessTTL, refreshTTL time.Duration, users user.Service) Service {
	return &ServiceImpl{secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL, users: users}
}

type TokenPair struct {
	UserID       uuid.UUID
	AccessToken  string
	RefreshToken string
	RefreshJTI   string
}

func (s *ServiceImpl) Login(dto dto.LoginBody) (*TokenPair, error) {
	u, err := s.users.GetByEmail(dto.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := s.users.VerifyPassword(u.Password, dto.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}
	jti := uuid.NewString()
	at, rt, err := utils.GenerateTokens(u.ID, s.secret, s.accessTTL, s.refreshTTL, jti)
	if err != nil {
		return nil, err
	}
	return &TokenPair{UserID: u.ID, AccessToken: at, RefreshToken: rt, RefreshJTI: jti}, nil
}

func (s *ServiceImpl) Refresh(refreshAuthHeader string) (*TokenPair, error) {
	token := refreshAuthHeader
	if len(token) > 7 && (strings.HasPrefix(token, "Bearer ") || strings.HasPrefix(token, "bearer ")) {
		token = token[7:]
	}
	claims, err := utils.ParseToken(token, s.secret)
	if err != nil {
		return nil, err
	}
	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return nil, jwt.ErrTokenExpired
	}
	newJTI := uuid.NewString()
	at, rt, err := utils.GenerateTokens(claims.UserID, s.secret, s.accessTTL, s.refreshTTL, newJTI)
	if err != nil {
		return nil, err
	}
	return &TokenPair{UserID: claims.UserID, AccessToken: at, RefreshToken: rt, RefreshJTI: newJTI}, nil
}
