package graph

import "github.com/reginbald/gqlgen-dataloader-subscription/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Todos map[string]*model.Todo
	Users map[string]*model.User
}
