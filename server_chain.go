package grpcx

import (
	"context"

	"google.golang.org/grpc"
)

// ServerUnaryChain build the multi server unary interceptors into one interceptor chain.
func ServerUnaryChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = buildServerUnaryInterceptor(interceptors[i], info, chain)
		}
		return chain(ctx, req)
	}
}

// build is the interceptor chain helper
func buildServerUnaryInterceptor(c grpc.UnaryServerInterceptor, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) grpc.UnaryHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return c(ctx, req, info, handler)
	}
}

// WithServerUnaryInterceptors is a grpc.Server config option that accepts multiple unary interceptors.
func WithServerUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.UnaryInterceptor(ServerUnaryChain(interceptors...))
}

// ServerStreamChain build the multi server stream interceptors into one interceptor chain.
func ServerStreamChain(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = buildServerStreamInterceptor(interceptors[i], info, chain)
		}
		return chain(srv, ss)
	}
}

func buildServerStreamInterceptor(c grpc.StreamServerInterceptor, info *grpc.StreamServerInfo, handler grpc.StreamHandler) grpc.StreamHandler {
	return func(srv interface{}, stream grpc.ServerStream) error {
		return c(srv, stream, info, handler)
	}
}

// WithServerStreamInterceptors is a grpc.Server config option that accepts multiple stream interceptors.
func WithServerStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.StreamInterceptor(ServerStreamChain(interceptors...))
}
