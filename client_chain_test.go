package grpcx

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestChainUnaryClient(t *testing.T) {
	ignoredMd := metadata.Pairs("foo", "bar")
	parentOpts := []grpc.CallOption{grpc.Header(&ignoredMd)}
	reqMessage := "request"
	replyMessage := "reply"
	outputError := fmt.Errorf("some error")

	first := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if ctx.Value("parent").(int) != someValue {
			t.Fatal("first interceptor must know the parent context value")
		}

		if !reflect.DeepEqual(someServiceName, method) {
			t.Fatal("first interceptor must know the someUnaryServerInfo")
		}

		if len(opts) != 1 {
			t.Fatal("first should see parent CallOptions")
		}

		ctx = context.WithValue(ctx, "first", 1)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
	second := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if ctx.Value("parent").(int) != someValue {
			t.Fatal("second interceptor must know the parent context value")
		}

		if ctx.Value("first").(int) != 1 {
			t.Fatal("second interceptor must know the first context value")
		}

		if !reflect.DeepEqual(someServiceName, method) {
			t.Fatal("second interceptor must know the someUnaryServerInfo")
		}

		if len(opts) != 1 {
			t.Fatal("second should see parent CallOptions")
		}

		wrappedOpts := append(opts, grpc.FailFast(true))
		ctx = context.WithValue(ctx, "second", 1)
		return invoker(ctx, method, req, reply, cc, wrappedOpts...)
	}
	invoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		if ctx.Value("parent").(int) != someValue {
			t.Fatal("invoker interceptor must know the parent context value")
		}

		if ctx.Value("first").(int) != 1 {
			t.Fatal("invoker interceptor must know the first context value")
		}

		if ctx.Value("second").(int) != 1 {
			t.Fatal("invoker interceptor must know the second context value")
		}

		if !reflect.DeepEqual(someServiceName, method) {
			t.Fatal("invoker interceptor must know the someUnaryServerInfo")
		}

		if len(opts) != 2 {
			t.Fatal("invoker should see parent CallOptions")
		}

		return outputError
	}
	chain := ClientUnaryChain(first, second)
	err := chain(parentContext, someServiceName, reqMessage, replyMessage, nil, invoker, parentOpts...)
	if !reflect.DeepEqual(err, outputError) {
		t.Fatal("chain must return invokers's error")
	}
}

type fakeClientStream struct {
	grpc.ClientStream
}

func TestChainStreamClient(t *testing.T) {
	ignoredMd := metadata.Pairs("foo", "bar")
	parentOpts := []grpc.CallOption{grpc.Header(&ignoredMd)}
	clientStream := &fakeClientStream{}
	fakeStreamDesc := &grpc.StreamDesc{ClientStreams: true, ServerStreams: true, StreamName: someServiceName}

	first := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if ctx.Value("parent").(int) != someValue {
			t.Fatal("first interceptor must know the parent context value")
		}

		if !reflect.DeepEqual(someServiceName, method) {
			t.Fatal("first interceptor must know the someUnaryServerInfo")
		}

		if len(opts) != 1 {
			t.Fatal("first should see parent CallOptions")
		}

		ctx = context.WithValue(ctx, "first", 1)

		return streamer(ctx, desc, cc, method, opts...)
	}
	second := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if ctx.Value("parent").(int) != someValue {
			t.Fatal("second interceptor must know the parent context value")
		}

		if ctx.Value("first").(int) != 1 {
			t.Fatal("second interceptor must know the first context value")
		}

		if !reflect.DeepEqual(someServiceName, method) {
			t.Fatal("second interceptor must know the someUnaryServerInfo")
		}

		if len(opts) != 1 {
			t.Fatal("second should see parent CallOptions")
		}

		wrappedOpts := append(opts, grpc.FailFast(true))
		ctx = context.WithValue(ctx, "second", 1)
		return streamer(ctx, desc, cc, method, wrappedOpts...)
	}
	streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		if ctx.Value("parent").(int) != someValue {
			t.Fatal("invoker interceptor must know the parent context value")
		}

		if ctx.Value("first").(int) != 1 {
			t.Fatal("invoker interceptor must know the first context value")
		}

		if ctx.Value("second").(int) != 1 {
			t.Fatal("invoker interceptor must know the second context value")
		}

		if !reflect.DeepEqual(someServiceName, method) {
			t.Fatal("invoker interceptor must know the someUnaryServerInfo")
		}

		if len(opts) != 2 {
			t.Fatal("invoker should see parent CallOptions")
		}

		if !reflect.DeepEqual(fakeStreamDesc, desc) {
			t.Fatal("streamer must see the right StreamDesc")
		}

		return clientStream, nil
	}

	chain := ClientStreamChain(first, second)
	someStream, err := chain(parentContext, fakeStreamDesc, nil, someServiceName, streamer, parentOpts...)
	if err != nil {
		t.Fatal("chain must not return an error as nothing there reutrned it")
	}

	if !reflect.DeepEqual(someStream, clientStream) {
		t.Fatal("chain must return invokers's clientstream")
	}
}
