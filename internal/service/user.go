package service

import (
	"errors"
	"todo-api/internal/models"
	"todo-api/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(name, email, password string) (*models.UserResponse, error) {
	existingUser, _ := s.repo.GetByEmail(email)
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
	err = s.repo.Create(repoUser)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) GetUser(id string) (*models.UserResponse, error) {
	repoUser, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if repoUser == nil {
		return nil, errors.New("User not found")
	}

	return &models.UserResponse{
		ID:    repoUser.ID,
		Name:  repoUser.Name,
		Email: repoUser.Email,
	}, nil
}

func (s *UserService) GetAllUsers() ([]models.UserResponse, error) {
	repoUsers, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	userResponses := make([]models.UserResponse, len(repoUsers))
	for i, repoUser := range repoUsers {
		userResponses[i] = models.UserResponse{
			ID:    repoUser.ID,
			Name:  repoUser.Name,
			Email: repoUser.Email,
		}
	}

	return userResponses, nil
}

func (s *UserService) UpdateUser(id string, req models.UpdateUserRequest) (*models.UserResponse, error) {
	repoUser, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if repoUser == nil {
		return nil, errors.New("User not found")
	}

	user := models.ConvertFromRepositoryUser(*repoUser)

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Email != "" {
		if req.Email != user.Email {
			existingUser, _ := s.repo.GetByEmail(req.Email)
			if existingUser != nil {
				return nil, errors.New("Email already in use")
			}
		}
		user.Email = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	updatedRepoUser := user.ConvertToRepositoryUser()
	err = s.repo.Update(updatedRepoUser)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
