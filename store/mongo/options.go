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

type mid_sub string

// ListPrefix returns all keys that are prefixed with key
func MidSubPrefix(p string) store.ListOption {
	return func(o *store.ListOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, mid_sub(""), p)
		//l.Prefix = p
	}
}
