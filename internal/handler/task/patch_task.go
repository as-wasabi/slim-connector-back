package task

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"slim-connector-back/model"
	"time"
)

func (h *TaskHandler) PatchTask(c *gin.Context) {
	var task model.Task
	//taskID := c.Param("id") 将来的にurlからid取得 -> 今はurlからbindしてる。

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.UpdateAt = time.Now()

	filter := bson.D{{"_id", task.ID}}
	update := bson.D{{"$set", task}}

	result, err := h.collection.UpdateOne(context.Background(), filter, update)

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
