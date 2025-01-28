package task

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"slim-connector-back/internal"
)

type TaskHandler struct {
	collection *mongo.Collection
}

func NewTaskHandler(initializer *internal.Initializer) *TaskHandler {
	collection := initializer.Database.Collection("tasks")
	return &TaskHandler{collection: collection}
}

func (h *TaskHandler) InitRoute(group *gin.RouterGroup) {
	group.GET("/tasks", h.GetTask)
	group.POST("/tasks", h.CreateTask)

}
