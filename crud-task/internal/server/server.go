package server

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"

	"github.com/berezovskiydeval/crud-task/internal/delivery/rest"
	"github.com/berezovskiydeval/crud-task/internal/repository"
	"github.com/berezovskiydeval/crud-task/internal/service"
)

func NewServer(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	taskRepo := repository.NewTaskRepository(db)

	taskService := service.NewTaskService(taskRepo)

	handler := rest.NewHandler(taskService)

	r.POST("/tasks", handler.CreateTaskHandler)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
