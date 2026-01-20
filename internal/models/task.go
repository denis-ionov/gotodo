package models

type TaskStatus string

const (
	StatusNew        TaskStatus = "New"
	StatusInProgress TaskStatus = "In progress"
	StatusCompleted  TaskStatus = "Finished"
)

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	UserID      string     `json:"user_id,omitempty"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
}
