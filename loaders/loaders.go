package loaders

import (
	"context"
	"net/http"

	"github.com/reginbald/gqlgen-dataloader-subscription/graph/model"
	"github.com/reginbald/gqlgen-dataloader-subscription/repository"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	UserLoader *dataloadgen.Loader[string, *model.User]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(repo *repository.Repository) *Loaders {
	// define the data loader
	return &Loaders{
		UserLoader: dataloadgen.NewLoader(func(ctx context.Context, keys []string) ([]*model.User, []error) {
			res := make([]*model.User, 0, len(keys))
			errs := make([]error, 0, len(keys))

			for _, key := range keys {
				user, err := repo.GetUser(key)
				res = append(res, &model.User{
					ID:   user.ID.String(),
					Name: user.Name,
				})
				errs = append(errs, err)
			}
			return res, errs

		}),
	}
}

// Middleware injects data loaders into the context
func Middleware(repo *repository.Repository, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders(repo)
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
