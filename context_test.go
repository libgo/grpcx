package grpcx

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestGenContext(t *testing.T) {
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("auth", "abc")) // mock grpc incoming
	ctx = GenContext(ctx, "tid", "123")
	ctx = GenContext(ctx, "mode", "debug")

	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		t.Fatal("ctx must be outging")
	}

	if md.Get("auth")[0] != "abc" || md.Get("tid")[0] != "123" || md.Get("mode")[0] != "debug" {
		t.Fatal("ctx val missing")
	}
}
