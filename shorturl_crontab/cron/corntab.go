package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	"shorturl_crontab/data"
	"shorturl_crontab/pkg/db/mysql"
	"shorturl_crontab/pkg/db/redis"
	"shorturl_crontab/pkg/log"
	"time"
)

const DefaultUrlMapTTL = 30 * 86400

func Run() {
	setUrlMapMaxID()
	c := cron.New()
	c.AddFunc("0 3 * * *", setUrlMapMaxID)
	c.Run()
}

func setUrlMapMaxID() {
	tables := []string{"url_map", "url_map_user"}
	redisPool := redis.GetPool()
	client := redisPool.Get()
	defer redisPool.Put(client)

	d := data.NewData(mysql.GetDB())
	for _, t := range tables {
		id, err := d.GetMaxID(t)
		if err != nil {
			log.Error(err)
			continue
		}
		key := redis.GetKey(t, "max_id")
		err = client.SetEx(context.Background(), key, id, time.Second*time.Duration(DefaultUrlMapTTL)).Err()
		if err != nil {
			log.Error(err)
			continue
		}
	}
}
