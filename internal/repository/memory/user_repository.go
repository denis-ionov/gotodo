package memory

import (
	"todo-api/internal/repository"
)

type userRepository struct {
	users map[string]repository.User
}

func NewUserRepository() repository.UserRepository {
	return &userRepository{
		users: make(map[string]repository.User),
	}
}

func (r *userRepository) Create(user repository.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *userRepository) GetByID(id string) (*repository.User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*repository.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, nil
}

func (r *userRepository) GetAll() ([]repository.User, error) {
	users := make([]repository.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Update(user repository.User) error {
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
