package grpcx

import (
	"context"

	"google.golang.org/grpc"
)

// ClientUnaryChain build the multi client unary interceptors into one interceptor chain.
func ClientUnaryChain(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		chain := invoker
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = buildClientUnaryInterceptor(interceptors[i], chain)
		}
		return chain(ctx, method, req, reply, cc, opts...)
	}
}

func buildClientUnaryInterceptor(c grpc.UnaryClientInterceptor, invoker grpc.UnaryInvoker) grpc.UnaryInvoker {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return c(ctx, method, req, reply, cc, invoker, opts...)
	}
}

// WithClientUnaryInterceptors is a grpc.Client dial option that accepts multiple unary interceptors.
func WithClientUnaryInterceptors(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithUnaryInterceptor(ClientUnaryChain(interceptors...))
}

//  ClientStreamChain build the multi client stream interceptors into one interceptor chain.
func ClientStreamChain(interceptors ...grpc.StreamClientInterceptor) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		chain := streamer
		for i := len(interceptors) - 1; i >= 0; i-- {
			chain = buildClientStreamInterceptor(interceptors[i], chain)
		}
		return chain(ctx, desc, cc, method, opts...)
	}
}

func buildClientStreamInterceptor(c grpc.StreamClientInterceptor, streamer grpc.Streamer) grpc.Streamer {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return c(ctx, desc, cc, method, streamer, opts...)
	}
}

// WithClientStreamInterceptors is a grpc.Client dial option that accepts multiple stream interceptors.
func WithClientStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithStreamInterceptor(ClientStreamChain(interceptors...))
}
