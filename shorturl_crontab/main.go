package main

import (
	"flag"
	"shorturl_crontab/cron"
	"shorturl_crontab/pkg/config"
	"shorturl_crontab/pkg/db/mysql"
	"shorturl_crontab/pkg/db/redis"
	"shorturl_crontab/pkg/log"
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

	//初始化mysql
	mysql.InitMysql(cnf)
	//初始化redis
	redis.InitRedisPool(cnf)

	cron.Run()
}
