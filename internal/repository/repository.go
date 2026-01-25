package repository

type TaskRepository interface {
	Create(task Task) error
	GetByID(id string) (*Task, error)
	GetAll() ([]Task, error)
	GetByUserID(userID string) ([]Task, error)
	Update(task Task) error
	Delete(id string) error
}

type UserRepository interface {
	Create(user User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAll() ([]User, error)
	Update(user User) error
	Delete(id string) error
}

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	UserID      string `json:"user_id"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type Repository struct {
	Task TaskRepository
	User UserRepository
}
