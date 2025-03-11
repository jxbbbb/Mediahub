package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"shorturl/pkg/config"
	"shorturl/pkg/db/mysql"
	"shorturl/pkg/db/redis"
	"shorturl/pkg/log"
	"shorturl/proto"
	"shorturl/shorturl-server/cache"
	"shorturl/shorturl-server/data"
	"shorturl/shorturl-server/interceptor"
	"shorturl/shorturl-server/server"
)

var (
	configFile = flag.String("config", "dev.config.yaml", "")
)

func main() {
	flag.Parse()
	//初始化配置文件
	config.InitConfig(*configFile)
	cnf := config.GetConfig()

	log.SetLevel(cnf.Log.Level)
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetPrintCaller(true)

	logger := log.NewLogger()
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)
	//初始化mysql
	mysql.InitMysql(cnf)
	urlMapDataFactory := data.NewUrlMapDataFactory(logger, mysql.GetDB())
	//初始化redis
	redis.InitRedisPool(cnf)
	kvCacheFactory := cache.NewRedisCacheFactory(redis.GetPool())

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cnf.Server.IP, cnf.Server.Port))
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(interceptor.UnaryAuthInterceptor), grpc.StreamInterceptor(interceptor.StreamAuthInterceptor))
	service := server.NewService(cnf, logger, urlMapDataFactory, kvCacheFactory)
	proto.RegisterShortUrlServer(s, service)

	healthCheckSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthCheckSrv)

	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
