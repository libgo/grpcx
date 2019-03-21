package grpcx

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type validator interface {
	Validate() error
}

// Validator checks if unary incoming proto is Valid
func Validator(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
	if v, ok := req.(validator); ok {
		if err = v.Validate(); err != nil {
			return nil, errorConvertor(err, codes.InvalidArgument)
		}
	}
	return handler(ctx, req)
}

// ValidatorStream checks if stream incoming proto is Valid
func ValidatorStream(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	wrapper := &recvWrapper{ss}
	return handler(srv, wrapper)
}

type recvWrapper struct {
	grpc.ServerStream
}

func (s *recvWrapper) RecvMsg(m interface{}) (err error) {
	if err = s.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	if v, ok := m.(validator); ok {
		if err = v.Validate(); err != nil {
			return errorConvertor(err, codes.InvalidArgument)
		}
	}
	return nil
}
