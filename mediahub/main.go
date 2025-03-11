package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"mediahub/Controller"
	"mediahub/middleware"
	"mediahub/pkg/config"
	"mediahub/pkg/log"
	"mediahub/pkg/storage/cos"
	"mediahub/routers"
	"net/http"
)

var configFile = flag.String("config", "dev.config.yaml", "Path to config file")

func main() {
	flag.Parse()
	//初始化配置文件
	config.InitConfig(*configFile)
	cnf := config.GetConfig()
	fmt.Println("CDNDomain from config:", cnf.Cos.CDNDomain) // 添加调试信息

	log.SetLevel(cnf.Log.Level)
	log.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	log.SetPrintCaller(true)

	logger := log.NewLogger()
	logger.SetOutput(log.GetRotateWriter(cnf.Log.LogPath))
	logger.SetLevel(cnf.Log.Level)
	logger.SetPrintCaller(true)
	sf := cos.NewCosStorageFactory(cnf.Cos.BucketUrl, cnf.Cos.SecretId, cnf.Cos.SecretKey, cnf.Cos.CDNDomain)
	controller := Controller.NewController(sf, logger, cnf)

	gin.SetMode(cnf.Http.Mode)
	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/health", func(c *gin.Context) {})
	api := r.Group("/api")
	routers.InitRouters(api, controller)

	fs := http.FileServer(http.Dir("www"))
	r.NoRoute(func(c *gin.Context) {
		fs.ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "www/index.html")
	})
	r.Run(fmt.Sprintf("%s:%d", cnf.Http.IP, cnf.Http.Port))
}
