package handler

import (
	"slim-connector-back/internal"
	"slim-connector-back/internal/handler/OpenAI"
	"slim-connector-back/internal/handler/task"
	"slim-connector-back/internal/handler/user"
)

func InitRoute(initializer *internal.Initializer) []internal.Route {
	taskHandler := task.NewTaskHandler(initializer)
	return []internal.Route{
		user.NewUserHandler(initializer),
		taskHandler,
		OpenAI.NewOpenAIHandler(initializer, taskHandler),
	}
}
