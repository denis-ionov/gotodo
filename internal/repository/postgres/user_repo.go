package postgres

import (
	"database/sql"
	"todo-api/internal/repository"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user repository.User) error {
	query := `
		INSERT INTO users (id, name, email, password)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	)

	return err
}

func (r *userRepository) GetByID(id string) (*repository.User, error) {
	query := `
		SELECT id, name, email, password
		FROM users
		WHERE id = $1
	`

	var user repository.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*repository.User, error) {
	query := `
		SELECT id, name, email, password
		FROM users
		WHERE email = $1
	`

	var user repository.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAll() ([]repository.User, error) {
	query := `
		SELECT id, name, email, password
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []repository.User
	for rows.Next() {
		var user repository.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Update(user repository.User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3, password = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
	)

	return err
}

func (r *userRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
