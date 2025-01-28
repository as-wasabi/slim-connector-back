package handler

import (
	"slim-connector-back/internal"
	"slim-connector-back/internal/handler/task"
	"slim-connector-back/internal/handler/user"
)

func InitRoute(initializer *internal.Initializer) []internal.Route {
	return []internal.Route{
		user.NewUserHandler(initializer),
		task.NewTaskHandler(initializer),
	}
}
