package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	model "github.com/thangchung/go-coffeeshop/internal/auth/domain"
	repository "github.com/thangchung/go-coffeeshop/internal/auth/repo"
)

var (
	ErrEmailTaken    = errors.New("email already registered")
	ErrUsernameTaken = errors.New("username already taken")
	ErrInvalidCreds  = errors.New("invalid email or password")
	ErrUserNotFound  = errors.New("user not found")
)

// AuthService contains the core business logic.
type AuthService struct {
	repo        *repository.UserRepository
	jwtSecret   []byte
	expireHours int
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string, expireHours int) *AuthService {
	return &AuthService{
		repo:        repo,
		jwtSecret:   []byte(jwtSecret),
		expireHours: expireHours,
	}
}

// Register creates a new user account and returns a JWT token.
func (s *AuthService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Check if email is already taken
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("register check email: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	// Hash the password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Persist the new user
	user, err := s.repo.Create(req.Username, req.Email, string(hash))
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("new user registered")

	// Issue JWT
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

// Login verifies credentials and returns a JWT token.
func (s *AuthService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("login find user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCreds
	}

	// Compare password with stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	logrus.WithFields(logrus.Fields{
		"user_id":  user.ID,
		"username": user.Username,
	}).Info("user logged in")

	// Issue JWT
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{Token: token, User: *user}, nil
}

// generateToken creates a signed JWT for the given user.
func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Duration(s.expireHours) * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}
