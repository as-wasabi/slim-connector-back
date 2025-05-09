package task

import (
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"slim-connector-back/internal"
)

type TaskHandler struct {
	collection *mongo.Collection
	node       *snowflake.Node
}

func NewTaskHandler(initializer *internal.Initializer) *TaskHandler {
	collection := initializer.Database.Collection("tasks")
	node, err := snowflake.NewNode(5)
	if err != nil {
		log.Fatalf("Failed to initialize Snowflake node: %v", err)
	}
	return &TaskHandler{
		collection: collection,
		node:       node,
	}
}

func (h *TaskHandler) InitRoute(group *gin.RouterGroup) {
	group.GET("/tasks", h.GetTask)
	group.GET("/tasks/analysis", h.TaskAnalysis)
	group.POST("/tasks", h.CreateTask)
	group.POST("/tasks/:id/create", h.CreateTaskFromDependence)
	group.PATCH("/tasks", h.PatchTask)

}
