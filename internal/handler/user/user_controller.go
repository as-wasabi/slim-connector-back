package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"slim-connector-back/internal"
	"slim-connector-back/model"
	"time"
)

type UserHandler struct {
	collection *mongo.Collection
}

func NewUserHandler(initializer *internal.Initializer) *UserHandler {
	collection := initializer.Database.Collection("users")
	return &UserHandler{collection: collection}
}
func (receiver UserHandler) InitRoute(group *gin.RouterGroup) {
	group.POST("/users", receiver.CreateUser)
	group.GET("/users", receiver.GetUsers)
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
