package repository

import (
	"todo-api/internal/models"
)

type TaskRepository interface {
	Create(task models.Task) error
	GetByID(id string) (*models.Task, error)
	GetAll() ([]models.Task, error)
	Update(task models.Task) error
	Delete(id string) error
}

type taskRepository struct {
	tasks map[string]models.Task
}

func NewTaskRepository() TaskRepository {
	return &taskRepository{
		tasks: make(map[string]models.Task),
	}
}

func (r *taskRepository) Create(task models.Task) error {
	r.tasks[task.ID] = task
	return nil
}

func (r *taskRepository) GetByID(id string) (*models.Task, error) {
	task, exists := r.tasks[id]
	if !exists {
		return nil, nil
	}

	return &task, nil
}

func (r *taskRepository) GetAll() ([]models.Task, error) {
	tasks := make([]models.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) Update(task models.Task) error {
	if _, exists := r.tasks[task.ID]; !exists {
		return nil
	}

	r.tasks[task.ID] = task
	return nil
}

func (r *taskRepository) Delete(id string) error {
	delete(r.tasks, id)
	return nil
}
