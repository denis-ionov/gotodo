package service

import (
	"errors"
	"todo-api/internal/models"
	"todo-api/internal/repository"

	"github.com/google/uuid"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(title, description, userID string) (*models.Task, error) {
	if title == "" {
		return nil, errors.New("Task title is required")
	}

	task := models.Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Status:      models.StatusNew,
		UserID:      userID,
	}

	err := s.repo.Create(task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) GetTask(id string) (*models.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errors.New("Task not found")
	}

	return task, nil
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	return s.repo.GetAll()
}

func (s *TaskService) UpdateTask(id string, req models.UpdateTaskRequest, userID string) (*models.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errors.New("Task not found")
	}

	if task.UserID != userID {
		return nil, errors.New("Access denied")
	}

	if req.Title != "" {
		task.Title = req.Title
	}

	if req.Description != "" {
		task.Description = req.Description
	}

	if req.Status != "" {
		task.Status = req.Status
	}

	err = s.repo.Update(*task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(id, userID string) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if task == nil {
		return errors.New("Task not found")
	}

	if task.UserID != userID {
		return errors.New("Access denied")
	}

	return s.repo.Delete(id)
}
