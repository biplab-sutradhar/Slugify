package services

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/biplab-sutradhar/slugify/api/internal/auth"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/gin-gonic/gin"
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

func MintAPIKey(svc *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not identified"})
			return
		}
		var req struct {
			Name string `json:"name"`
		}
		_ = c.ShouldBindJSON(&req)
		if req.Name == "" {
			req.Name = "Default"
		}
		key, err := svc.MintAPIKey(c, userID, req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"api_key": key})
	}
}
