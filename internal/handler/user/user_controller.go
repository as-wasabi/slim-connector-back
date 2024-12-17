package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"slim-connector-back/internal/repository"
	"slim-connector-back/model"
	"time"
)

type UserHandler struct {
	collection *mongo.Collection
}

func NewUserHandler(mongoDB *repository.MongoDB) *UserHandler {
	collection := mongoDB.Database.Collection("users")
	return &UserHandler{collection: collection}
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	var user model.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser := h.collection.FindOne(context.Background(), bson.M{"email": user.Email})

	if existingUser.Err() == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	user.CreatedAt = time.Now()

	result, err := h.collection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User Created Successfully",
		"user_id": result.InsertedID,
	})

}

func (h *UserHandler) GetUsers(c *gin.Context) {
	cursor, err := h.collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, context.Background())

	var users []model.User

	if err = cursor.All(context.Background(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, users)

}
