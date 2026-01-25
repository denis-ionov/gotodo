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
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(name, email, password string) (*models.UserResponse, error) {
	existingUser, _ := s.userRepo.GetByEmail(email)
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

	repoUser := user.ConvertToRepositoryUser()
	err = s.userRepo.Create(repoUser)
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
	repoUser, err := s.userRepo.GetByEmail(email)
	if err != nil || repoUser == nil {
		return "", errors.New("Wrong email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(repoUser.Password), []byte(password))
	if err != nil {
		return "", errors.New("Wrong email or password")
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = repoUser.ID
	claims["email"] = repoUser.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) GetUserByID(userID string) (*models.UserResponse, error) {
	repoUser, err := s.userRepo.GetByID(userID)
	if err != nil || repoUser == nil {
		return nil, errors.New("User not found")
	}

	return &models.UserResponse{
		ID:    repoUser.ID,
		Name:  repoUser.Name,
		Email: repoUser.Email,
	}, nil
}
