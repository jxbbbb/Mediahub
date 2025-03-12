package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"shorturl-proxy/pkg/config"
	"shorturl-proxy/pkg/log"
	"shorturl-proxy/proxy"
)

var configFile = flag.String("config", "dev.config.yaml", "")

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

	gin.SetMode(cnf.Http.Mode)
	r := gin.Default()
	r.GET("/health", func(*gin.Context) {})
	p := proxy.NewProxy(cnf, logger)
	public := r.Group("/p")
	public.GET("/:short_key", p.PublicProxy)
	user := r.Group("/u")
	user.GET("/:short_key", p.UserProxy)

	err := r.Run(fmt.Sprintf("%s:%d", cnf.Http.IP, cnf.Http.Port))
	if err != nil {
		log.Error(err)
	}
}
