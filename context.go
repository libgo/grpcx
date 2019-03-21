package grpcx

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// Context is a wrapper for context.Context
// TODO header and trailer should be propagable? [client <------ svc1 <---h/t--- svc2]
type Context struct {
	context.Context
	// rw
	// header
	// trailer
}

// GenContext always gen grpc outging context
func GenContext(ctx context.Context, kv ...string) *Context {
	if c, ok := ctx.(*Context); ok {
		if len(kv) != 0 {
			c.Context = metadata.AppendToOutgoingContext(c.Context, kv...)
		}
		return c
	}

	// incoming to outging
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		newMd := metadata.Pairs(kv...)

		// TODO maybe we should skip key with "grpc-" prefix
		for k, v := range md {
			newMd.Set(k, v...)
		}

		ctx = metadata.NewOutgoingContext(ctx, newMd)
		return &Context{
			Context: ctx,
		}
	}

	// outgoing kv append
	if _, ok := metadata.FromOutgoingContext(ctx); ok {
		if len(kv) != 0 {
			ctx = metadata.AppendToOutgoingContext(ctx, kv...)
		}
		return &Context{
			Context: ctx,
		}
	}

	return &Context{
		Context: metadata.NewOutgoingContext(ctx, metadata.Pairs(kv...)),
	}
}
