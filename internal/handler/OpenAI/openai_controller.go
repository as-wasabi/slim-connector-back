package OpenAI

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"slim-connector-back/internal"
	"slim-connector-back/internal/handler/task"
)

type OpenAIHandler struct {
	collection  *mongo.Collection
	TaskHandler *task.TaskHandler
}

func NewOpenAIHandler(initializer *internal.Initializer, taskHandler *task.TaskHandler) *OpenAIHandler {
	collection := initializer.Database.Collection("openai")
	return &OpenAIHandler{
		collection:  collection,
		TaskHandler: taskHandler,
	}
}

func (h *OpenAIHandler) InitRoute(group *gin.RouterGroup) {
	group.POST("/openai", h.ExtractedTask)
}
