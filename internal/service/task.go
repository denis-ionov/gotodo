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

	repoTask := task.ConvertToRepositoryTask()
	err := s.repo.Create(repoTask)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) GetTask(id string) (*models.Task, error) {
	repoTask, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if repoTask == nil {
		return nil, errors.New("Task not found")
	}

	task := models.ConvertFromRepositoryTask(*repoTask)
	return &task, nil
}

func (s *TaskService) GetAllTasks() ([]models.Task, error) {
	repoTasks, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	tasks := make([]models.Task, len(repoTasks))
	for i, repoTask := range repoTasks {
		tasks[i] = models.ConvertFromRepositoryTask(repoTask)
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(id string, req models.UpdateTaskRequest, userID string) (*models.Task, error) {
	repoTask, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if repoTask == nil {
		return nil, errors.New("Task not found")
	}

	if repoTask.UserID != userID {
		return nil, errors.New("Access denied")
	}

	task := models.ConvertFromRepositoryTask(*repoTask)

	if req.Title != "" {
		task.Title = req.Title
	}

	if req.Description != "" {
		task.Description = req.Description
	}

	if req.Status != "" {
		task.Status = req.Status
	}

	updatedRepoTask := task.ConvertToRepositoryTask()
	err = s.repo.Update(updatedRepoTask)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *TaskService) DeleteTask(id, userID string) error {
	repoTask, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if repoTask == nil {
		return errors.New("Task not found")
	}

	if repoTask.UserID != userID {
		return errors.New("Access denied")
	}

	return s.repo.Delete(id)
}
