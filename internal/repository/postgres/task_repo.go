package postgres

import (
	"database/sql"
	"todo-api/internal/repository"
)

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) repository.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task repository.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, user_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.UserID,
	)

	return err
}

func (r *taskRepository) GetByID(id string) (*repository.Task, error) {
	query := `
		SELECT id, title, description, status, user_id
		FROM tasks
		WHERE id = $1
	`

	var task repository.Task
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.UserID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *taskRepository) GetAll() ([]repository.Task, error) {
	query := `
		SELECT id, title, description, status, user_id
		FROM tasks
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []repository.Task
	for rows.Next() {
		var task repository.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.UserID,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) GetByUserID(userID string) ([]repository.Task, error) {
	query := `
		SELECT id, title, description, status, user_id
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []repository.Task
	for rows.Next() {
		var task repository.Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.UserID,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) Update(task repository.Task) error {
	query := `
		UPDATE tasks
		SET title = $2, description = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
	)

	return err
}

func (r *taskRepository) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
