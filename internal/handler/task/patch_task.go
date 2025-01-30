package task

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"slim-connector-back/model"
	"time"
)

func (h *TaskHandler) UpdateTask(c *gin.Context, task *model.Task) (*mongo.UpdateResult, error) {
	taskID := c.Param("id")
	task.UpdateAt = time.Now()

	filter := bson.D{{"_id", taskID}}
	update := bson.D{{"$set", task}}

	result, err := h.collection.UpdateOne(context.Background(), filter, update)

	return result, err
}

func (h *TaskHandler) PatchTask(c *gin.Context) {
	task, err := BindTaskJson(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.UpdateTask(c, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}
