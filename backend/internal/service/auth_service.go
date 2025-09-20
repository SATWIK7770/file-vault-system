package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
    ErrUserExists   = errors.New("user already exists")
    ErrEmailExists  = errors.New("email already exists")
    ErrInvalidCreds = errors.New("invalid username or password")
)


type AuthService struct {
	userRepo *repository.UserRepository
	jwtKey   []byte
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	return &AuthService{userRepo: userRepo, jwtKey: []byte(secret)}
}

// SignUp creates a new user
func (s *AuthService) SignUp(username, email, password string) error {
    // Check if username exists
    existingUser, _ := s.userRepo.GetByUsername(username)
    if existingUser != nil {
        return ErrUserExists
    }

    // Check if email exists
    existingEmail, _ := s.userRepo.GetByEmail(email)
    if existingEmail != nil {
        return ErrEmailExists
    }

    // Hash password
    hashed, err := models.HashPassword(password)
    if err != nil {
        return err
    }

    // Create user
    user := &models.User{
        Username:     username,
        Email:        email,
        Password: 	  hashed,
    }

    return s.userRepo.Create(user)
}


func (s *AuthService) SignUpAndGenerateToken(username, email, password string) (uint, string, error) {
    // Reuse your signup logic
    err := s.SignUp(username, email, password)
    if err != nil {
        return 0, "", err
    }

    // Fetch newly created user
    user, err := s.userRepo.GetByUsername(username)
    if err != nil || user == nil {
        return 0, "", errors.New("failed to fetch newly created user")
    }

    // Create JWT
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString(s.jwtKey)
    if err != nil {
        return 0, "", err
    }

    return user.ID, signed, nil
}


func (s *AuthService) SignIn(username, password string) (string, uint, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil || user == nil {
		return "", 0, ErrInvalidCreds
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", 0, ErrInvalidCreds
	}

	// Create JWT
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", 0, err
	}

	return signed, user.ID, nil
}


func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// JWT PARSE
func ParseToken(tokenStr string, jwtKey []byte) (uint, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidCreds
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid token payload")
	}

	return uint(userID), nil
}
