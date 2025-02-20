package OpenAI

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"slim-connector-back/internal"
)

type OpenAIHandler struct {
	collection *mongo.Collection
}

func NewOpenAIHandler(initializer *internal.Initializer) *OpenAIHandler {
	collection := initializer.Database.Collection("openai")
	return &OpenAIHandler{collection: collection}
}

func (h *OpenAIHandler) InitRoute(group *gin.RouterGroup) {
	group.GET("/openai", h.ExtractedTask)
}
