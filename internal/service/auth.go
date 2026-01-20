package service

import (
	"errors"
	"time"
	"todo-api/internal/models"
	"todo-api/internal/repository"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository repository.UserRepository
	jwtSecret      string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepository: userRepo,
		jwtSecret:      jwtSecret,
	}
}

func (s *AuthService) Register(name, email, password string) (*models.UserResponse, error) {
	existingUser, _ := s.userRepository.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("User with this email exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepository.GetByEmail(email)
	if err != nil || user == nil {
		return "", errors.New("Wrong email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("Wrong email or password")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) GetUserByID(userID string) (*models.UserResponse, error) {
	user, err := s.userRepository.GetByID(userID)
	if err != nil || user == nil {
		return nil, errors.New("User not found")
	}

	return &models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
