package repository

import (
	"gorm.io/gorm"

	"github.com/berezovskiydeval/crud-task/internal/domain"
)

type TaskRepository interface {
	CreateTask(task *domain.Task) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(task *domain.Task) error {
	return r.db.Create(task).Error
}
