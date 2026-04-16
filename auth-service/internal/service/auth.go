package service

import (
	"context"
	"errors"

	"github.com/kekus228swaga/orderflow/auth-service/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo user.Repository
}

func NewAuthService(userRepo user.Repository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, req user.RegisterRequest) (*user.User, error) {
	// Проверка: существует ли пользователь
	existing, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаём пользователя
	return s.userRepo.Create(ctx, req.Email, string(hash))
}

func (s *AuthService) Login(ctx context.Context, req user.LoginRequest) (*user.User, error) {
	u, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("invalid credentials")
	}

	// Сравниваем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return u, nil
}
