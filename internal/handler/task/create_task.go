package task

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	task.ID = h.node.Generate().String()
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

func (h *TaskHandler) CreateTaskFromDependence(c *gin.Context) {
	// Paramに変更できるはず～
	ParentID := c.Param("ParentID")
	task, err := BindTaskJson(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	task.Parent = ParentID

	if h.node == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Snowflake node is not initialized"})
		return
	}

	childTaskID := h.node.Generate().String()

	childTask := model.Task{
		ID:        childTaskID,
		Parent:    ParentID,
		Context:   task.Context,
		CreatedAt: time.Now(),
		Start:     task.Start,
		End:       task.End,
	}

	_, err = h.collection.InsertOne(context.Background(), childTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create child task"})
		return
	}

	update := bson.M{"$push": bson.M{"children": childTaskID}}
	_, err = h.collection.UpdateOne(context.Background(), bson.M{"id": ParentID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update parent task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Child task created successfully",
		"task_id":   childTaskID,
		"parent_id": ParentID,
	})

}
func (h *TaskHandler) CreateTaskFromAIResponse(task *model.Task) error {
	task.CreatedAt = time.Now()
	if h.node == nil {
		return fmt.Errorf("snowflake node is not initialized")
	}
	task.ID = h.node.Generate().String()
	_, err := h.collection.InsertOne(context.Background(), task)
	return err
}
