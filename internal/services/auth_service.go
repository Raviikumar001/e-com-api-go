// internal/services/auth_service.go

package services

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/Raviikumar001/e-com-api-go/internal/database"
	"github.com/Raviikumar001/e-com-api-go/internal/repositories"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    repositories.UserRepository
	jwtService  *JWTService
	userService *UserService
}

func NewAuthService(userRepo repositories.UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{userRepo: userRepo, jwtService: jwtService}
}

func (as *AuthService) Login(c *fiber.Ctx, username string, password string) (string, error) {
	// Find user by username
	user, err := as.userService.FindByUsername(c, username)
	if err != nil {
		return "", err
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fiber.ErrUnauthorized
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) Authenticate(username string, password string) (string, error) {
	var user database.User
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", err
	}
	if !user.ComparePassword(password) {
		return "", errors.New("invalid password")
	}
	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) ValidateToken(c *fiber.Ctx) (*JWTClaims, error) {
	auth := c.Get("Authorization")
	if auth == "" {
		return nil, errors.New("no authorization header")
	}

	// Check if the header starts with "Bearer "
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return nil, errors.New("invalid authorization header format")
	}

	// Extract the token by removing the "Bearer " prefix
	tokenString := auth[len(prefix):]
	return s.jwtService.ValidateToken(tokenString)
}
