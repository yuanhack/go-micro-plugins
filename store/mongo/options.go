package mongo

import (
	"context"
	"github.com/macheal/go-micro/v2/store"
)

type uri string

// Token sets the cloudflare api token
func URI(t string) store.Option {
	return func(o *store.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, uri(""), t)
	}
}
