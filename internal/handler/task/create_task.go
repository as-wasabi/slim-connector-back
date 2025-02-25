package task

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"slim-connector-back/model"
	"time"
)

func BindTaskJson(c *gin.Context) (*model.Task, error) {
	var task model.Task
	if err := c.BindJSON(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (h *TaskHandler) InsertTaskData(task *model.Task) (*mongo.InsertOneResult, error) {
	task.CreatedAt = time.Now()
	result, err := h.collection.InsertOne(context.Background(), task)

	return result, err
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	task, err := BindTaskJson(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 既存タスクチェックする？id生成的にダブんなくね???
	// existingTask := h.collection.FindOne(context.Background(), bson.M{"id": task.ID})

	result, err := h.InsertTaskData(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Task Created Successfully",
		"user_id": result.InsertedID,
	})
}
