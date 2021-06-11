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

// ListSubstr returns all keys that are prefixed with key
func SetListSubstr(p string) store.ListOption {
	return func(o *store.ListOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "ListSubstr", p)
	}
}

func GetListSubstr(ctx context.Context) string {
	var s string
	if v := ctx.Value("ListSubstr"); v != nil {
		s = v.(string)
	}
	return s
}

// ReadSubstr returns all keys that are prefixed with key
func SetReadSubstr(p string) store.ReadOption {
	return func(o *store.ReadOptions) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "ReadSubstr", p)
	}
}

func GetReadSubstr(ctx context.Context) string {
	var s string
	if v := ctx.Value("ReadSubstr"); v != nil {
		s = v.(string)
	}
	return s
}
