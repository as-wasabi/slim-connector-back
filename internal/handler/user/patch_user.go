package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"slim-connector-back/model"
	"time"
)

func (h *UserHandler) PatchUser(c *gin.Context) {
	var user model.User
	//userID := c.Param("id")

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.UpdateAt = time.Now()

	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", user}}

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
