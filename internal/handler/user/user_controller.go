package user

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"slim-connector-back/internal"
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
