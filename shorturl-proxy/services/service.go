package services

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func AppendBearerToken(ctx context.Context, accessToken string) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+accessToken)
	return metadata.NewOutgoingContext(ctx, md)
}
