package repository

import (
	"todo-api/internal/models"
)

type UserRepository interface {
	Create(user models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user models.User) error
	Delete(id string) error
}

type userRepository struct {
	users map[string]models.User
}

func NewUserRepository() UserRepository {
	return &userRepository{
		users: make(map[string]models.User),
	}
}

func (r *userRepository) Create(user models.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, nil
}

func (r *userRepository) GetAll() ([]models.User, error) {
	users := make([]models.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Update(user models.User) error {
	if _, exists := r.users[user.ID]; !exists {
		return nil
	}

	r.users[user.ID] = user
	return nil
}

func (r *userRepository) Delete(id string) error {
	delete(r.users, id)
	return nil
}
