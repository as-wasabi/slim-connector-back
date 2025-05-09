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

func BindUserJson(c *gin.Context) (*model.User, error) {
	var user model.User

	if err := c.BindJSON(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (h *UserHandler) ExistingUserCheck(user *model.User) error {
	existingUser := h.collection.FindOne(context.Background(), bson.M{"email": user.Email})

	if existingUser.Err() == nil {
		return existingUser.Err()
	}
	return nil
}

func (h *UserHandler) InsertUserdata(user *model.User) (*mongo.InsertOneResult, error) {
	user.CreatedAt = time.Now()

	result, err := h.collection.InsertOne(context.Background(), user)
	return result, err
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	user, err := BindUserJson(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = h.ExistingUserCheck(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	result, err := h.InsertUserdata(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// debug用
	c.JSON(http.StatusCreated, gin.H{
		"message": "User Created Successfully",
		"user_id": result.InsertedID,
	})

}
