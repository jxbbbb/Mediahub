package server

import (
	"context"
	"shorturl/pkg/config"
	"shorturl/pkg/log"
	"shorturl/pkg/utils"
	"shorturl/pkg/zerror"
	"shorturl/proto"
	"shorturl/shorturl-server/cache"
	"shorturl/shorturl-server/data"
	"time"
)

type shortURLService struct {
	proto.UnimplementedShortUrlServer
	config            *config.Config
	log               log.ILogger
	UrlMapDataFactory data.IUrlMapDataFactory
	KvCacheFactory    cache.CacheFactory
	//cache
	//db
}

func NewService(config *config.Config, log log.ILogger, urlMapDataFactory data.IUrlMapDataFactory, KvCacheFactory cache.CacheFactory) proto.ShortUrlServer {
	return &shortURLService{
		config:            config,
		log:               log,
		UrlMapDataFactory: urlMapDataFactory,
		KvCacheFactory:    KvCacheFactory,
	}
}

func (s *shortURLService) GetShortUrl(ctx context.Context, in *proto.Url) (*proto.Url, error) {
	isPublic := true
	if in.UserID != 0 {
		isPublic = false
	}
	if in.Url == "" {
		err := zerror.NewByMsg("参数为空")
		return nil, err
	}
	if !utils.IsUrl(in.Url) {
		err := zerror.NewByMsg("参数不合法")
		s.log.Error(err)
		return nil, err
	}
	//根据长链接查询数据库，记录是否存在（可以考虑用缓存，长链做key，或者短链做key，但功能主要是添加，查询业务较小，暂不考虑）
	data := s.UrlMapDataFactory.NewUrlMapData(isPublic)
	entity, err := data.GetByOriginal(in.Url)
	if err != nil {
		s.log.Error(zerror.NewByErr(err))
		return nil, err
	}
	if entity.ShortKey == "" {
		//新增记录
		id, err := data.GenerateID(in.GetUserID(), time.Now().Unix())
		if err != nil {
			s.log.Error(zerror.NewByErr(err))
			return nil, err
		}
		entity.ShortKey = utils.ToBase62(id)
		entity.OriginalUrl = in.Url
		entity.ID = id
		entity.UpdatedAt = time.Now().Unix()
		err = data.Update(entity)
		if err != nil {
			s.log.Error(zerror.NewByErr(err))
			return nil, err
		}
	}
	keyPrefix := ""
	domain := s.config.ShortDomain
	if !isPublic {
		keyPrefix = "user"
		domain = s.config.UserShortDomain
	}
	kvCache := s.KvCacheFactory.NewKVCache()
	defer kvCache.Destroy()
	key := keyPrefix + entity.ShortKey
	err = kvCache.Set(key, entity.OriginalUrl, cache.DefaultTTL)
	if err != nil {
		s.log.Error(zerror.NewByErr(err))
		return nil, err
	}
	return &proto.Url{
		Url:    domain + entity.ShortKey,
		UserID: in.UserID,
	}, nil
}
func (s *shortURLService) GetOriginUrl(ctx context.Context, in *proto.ShortKey) (*proto.Url, error) {
	isPublic := true
	if in.UserID != 0 {
		isPublic = false
	}
	if in.Key == "" {
		err := zerror.NewByMsg("参数为空")
		return nil, err
	}

	id := utils.ToBase10(in.Key)
	if id == 0 {
		err := zerror.NewByMsg("ID为空")
		s.log.Error(err)
		return nil, err
	}

	keyPrefix := ""
	if !isPublic {
		keyPrefix = "user_"
	}
	kvCache := s.KvCacheFactory.NewKVCache()
	defer kvCache.Destroy()
	key := keyPrefix + in.Key

	data := s.UrlMapDataFactory.NewUrlMapData(isPublic)
	originalUrl, err := kvCache.Get(key)
	if err != nil {
		s.log.Error(zerror.NewByErr(err))
		return nil, err
	}
	if originalUrl == "" {
		entity, err := data.GetByID(id)
		if err != nil {
			s.log.Error(zerror.NewByErr(err))
			return nil, err
		}
		originalUrl = entity.OriginalUrl
	}
	err = kvCache.Set(key, originalUrl, cache.DefaultTTL)
	if err != nil {
		s.log.Error(zerror.NewByErr(err))
		return nil, err
	}
	err = data.IncrementTimes(id, 1, time.Now().Unix())
	if err != nil {
		s.log.Warning(err)
		err = nil
	}
	return &proto.Url{
		Url:    originalUrl,
		UserID: in.UserID,
	}, nil
}
