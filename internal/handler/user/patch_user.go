package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"slim-connector-back/model"
	"time"
)

func (h *UserHandler) UpdateUser(user *model.User, c *gin.Context) (*mongo.UpdateResult, error) {
	user.UpdateAt = time.Now()
	userID := c.Param("id")
	filter := bson.D{{"_id", userID}}
	update := bson.D{{"$set", user}}

	result, err := h.collection.UpdateOne(context.Background(), filter, update)

	return result, err
}

func (h *UserHandler) PatchUser(c *gin.Context) {
	user, err := BindUserJson(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.UpdateUser(user, c)
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
