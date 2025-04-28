package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/berezovskiydeval/crud-task/internal/service"
)

var (
	tasksCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tasks_created_total",
		Help: "Total number of tasks created",
	})
	taskCreationDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "task_creation_duration_seconds",
		Help: "Duration of task creation in seconds",
	})
)

func init() {
	prometheus.MustRegister(tasksCreated)
	prometheus.MustRegister(taskCreationDuration)
}

type Handler struct {
	taskService service.TaskService
}

func NewHandler(taskService service.TaskService) *Handler {
	return &Handler{
		taskService: taskService,
	}
}

func (h *Handler) CreateTaskHandler(c *gin.Context) {
	startTime := time.Now()

	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.taskService.CreateTask(input.Title, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tasksCreated.Inc()
	taskCreationDuration.Observe(time.Since(startTime).Seconds())

	c.JSON(http.StatusCreated, task)
}
