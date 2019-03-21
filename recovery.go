package grpcx

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryLogFunc is default recovery log func, MUST be set before using Recovery interceptor.
var RecoveryLogFunc func(context.Context, interface{})

// Recovery used for first chain of grpc server unary interceptor for panic recover.
func Recovery(logFunc ...func(context.Context, interface{})) grpc.UnaryServerInterceptor {
	l := RecoveryLogFunc
	if len(logFunc) != 0 && logFunc[0] != nil {
		l = logFunc[0]
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				// if panic, set custom error to 'err', in order that client and sense it.
				err = status.Errorf(codes.Internal, "%s panic", info.FullMethod)
				if l != nil {
					l(ctx, req)
				}
			}
		}()
		return handler(ctx, req)
	}
}

// RecoveryStream used for first chain of grpc server stream interceptor for panic recover.
func RecoveryStream(logFunc ...func(context.Context, interface{})) grpc.StreamServerInterceptor {
	l := RecoveryLogFunc
	if len(logFunc) != 0 && logFunc[0] != nil {
		l = logFunc[0]
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "%s panic", info.FullMethod)
				if l != nil {
					l(ss.Context(), srv)
				}
			}
		}()
		return handler(srv, ss)
	}
}
