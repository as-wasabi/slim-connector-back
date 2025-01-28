package task

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"slim-connector-back/model"
	"time"
)

func (h *TaskHandler) CreateTask(c *gin.Context) {

	var task model.Task

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 既存タスクチェックする？id生成的にダブんなくね???
	// existingTask := h.collection.FindOne(context.Background(), bson.M{"id": task.ID})

	task.CreatedAt = time.Now()

	result, err := h.collection.InsertOne(context.Background(), task)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task Created Successfully",
		"user_id": result.InsertedID,
	})
}
