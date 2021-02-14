package mongo

import (
	"context"
	"github.com/micro/go-micro/v2/store"
)

type uri string

// Token sets the cloudflare api token
func URI(t string) store.Option {
	return func(o *store.Options) {
		o.Context = context.WithValue(o.Context, uri(""), t)
	}
}
