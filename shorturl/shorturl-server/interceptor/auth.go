package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"shorturl/pkg/config"
	"shorturl/pkg/zerror"
	"strings"
)

// TODO 此处通过拦截器实现鉴权
func UnaryAuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	//需要跳过健康检查
	if info.FullMethod != "/grpc.health.v1.Health/Check" {
		err = oauth2Vaild(ctx)
		if err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}

func StreamAuthInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := oauth2Vaild(ss.Context())
	if err != nil {
		return err
	}
	return handler(srv, ss)
}

func oauth2Vaild(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return zerror.NewByMsg("no metadata")
	}
	authorization := md["authorization"]
	if len(authorization) < 1 {
		return zerror.NewByMsg("no authorization")
	}
	token := strings.Trim(authorization[0], "Bearer ")
	cnf := config.GetConfig()
	if cnf.Server.AccessToken != token {
		return zerror.NewByMsg("身份验证失败")
	}
	return nil
}
