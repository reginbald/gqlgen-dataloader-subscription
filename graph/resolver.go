package graph

import (
	"github.com/reginbald/gqlgen-dataloader-subscription/repository"
)

type Resolver struct {
	Repo         *repository.Repository
	EventChannel <-chan int
}
