package grpcx

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Auth is helper interceptor for JWT or RBAC to check incoming ctx metadata.
func Auth(checkFunc func(context.Context) error) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		if checkFunc != nil {
			if err = checkFunc(ctx); err != nil {
				err = errorConvertor(err, codes.Unauthenticated)
				return
			}
		}
		return handler(ctx, req)
	}
}

// AuthStream is same as Auth but for server stream interceptor.
func AuthStream(checkFunc func(context.Context) error) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		if checkFunc != nil {
			if err = checkFunc(ss.Context()); err != nil {
				err = errorConvertor(err, codes.Unauthenticated)
				return
			}
		}
		return handler(srv, ss)
	}
}
