package memory

import (
	"todo-api/internal/repository"
)

type taskRepository struct {
	tasks map[string]repository.Task
}

func NewTaskRepository() repository.TaskRepository {
	return &taskRepository{
		tasks: make(map[string]repository.Task),
	}
}

func (r *taskRepository) Create(task repository.Task) error {
	r.tasks[task.ID] = task
	return nil
}

func (r *taskRepository) GetByID(id string) (*repository.Task, error) {
	task, exists := r.tasks[id]
	if !exists {
		return nil, nil
	}

	return &task, nil
}

func (r *taskRepository) GetAll() ([]repository.Task, error) {
	tasks := make([]repository.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *taskRepository) GetByUserID(userID string) ([]repository.Task, error) {
	var userTasks []repository.Task
	for _, task := range r.tasks {
		if task.UserID == userID {
			userTasks = append(userTasks, task)
		}
	}

	return userTasks, nil
}

func (r *taskRepository) Update(task repository.Task) error {
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
