package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/biplab-sutradhar/slugify/api/internal/auth"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService struct {
	users     db.UserRepository
	apiKeys   db.APIKeyRepository
	jwtSecret string
}

func NewAuthService(users db.UserRepository, apiKeys db.APIKeyRepository, jwtSecret string) *AuthService {
	return &AuthService{users: users, apiKeys: apiKeys, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (models.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if existing, _ := s.users.GetUserByEmail(ctx, email); existing.ID != "" {
		return models.AuthResponse{}, ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.AuthResponse{}, err
	}

	user := models.User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hash),
		Name:         strings.TrimSpace(req.Name),
		CreatedAt:    time.Now(),
	}
	if err := s.users.CreateUser(ctx, user); err != nil {
		return models.AuthResponse{}, err
	}

	apiKey, err := s.mintAPIKey(ctx, user.ID, "Default")
	if err != nil {
		return models.AuthResponse{}, err
	}

	token, err := auth.GenerateUserToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return models.AuthResponse{}, err
	}

	user.PasswordHash = ""
	return models.AuthResponse{Token: token, User: user, ApiKey: apiKey}, nil
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (models.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := s.users.GetUserByEmail(ctx, email)
	if err != nil {
		return models.AuthResponse{}, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return models.AuthResponse{}, ErrInvalidCredentials
	}

	token, err := auth.GenerateUserToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return models.AuthResponse{}, err
	}

	user.PasswordHash = ""
	return models.AuthResponse{Token: token, User: user}, nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (models.User, error) {
	u, err := s.users.GetUserByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}
	u.PasswordHash = ""
	return u, nil
}

func (s *AuthService) MintAPIKey(ctx context.Context, userID, name string) (string, error) {
	return s.mintAPIKey(ctx, userID, name)
}

func (s *AuthService) mintAPIKey(ctx context.Context, userID, name string) (string, error) {
	key, err := auth.GenerateAPIKey()
	if err != nil {
		return "", err
	}
	apiKey := models.APIKey{
		ID:        uuid.New().String(),
		UserID:    userID,
		Key:       key,
		Name:      name,
		Scope:     "default",
		IsActive:  true,
		CreatedAt: time.Now(),
		Usage:     0,
	}
	if err := s.apiKeys.CreateAPIKey(ctx, apiKey); err != nil {
		return "", err
	}
	return key, nil
}
