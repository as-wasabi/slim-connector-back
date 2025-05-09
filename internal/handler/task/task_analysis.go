package task

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"slim-connector-back/model"
)

func (h *TaskHandler) TaskAnalysis(c *gin.Context) {
	tasks := make(map[string]model.Task)
	progress := make(map[string]float64)

	cursor, err := h.collection.Find(context.Background(), bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"error": "Not Found cursor"})
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bson.M{"error": "よくわからんエラー"})
		}
	}(cursor, context.Background())

	var taskList []model.Task
	if err = cursor.All(context.Background(), &taskList); err != nil {
		c.JSON(http.StatusInternalServerError, bson.M{"error": "Not Found Tasks"})
	}

	for _, task := range taskList {
		tasks[task.ID] = task
	}

	for _, task := range taskList {
		if len(task.Children) == 0 {
			progress[task.ID] = task.CompleteRatio
		} else {
			totalRatio := 0.0

			for _, childID := range task.Children {
				if childTask, exists := tasks[childID]; exists {
					totalRatio += childTask.CompleteRatio
				}
			}
			progress[task.ID] = totalRatio / float64(len(task.Children))
		}
	}
	c.JSON(http.StatusOK, progress)
	//return progress, nil
}
