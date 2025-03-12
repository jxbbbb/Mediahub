package proxy

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"shorturl-proxy/pkg/config"
	"shorturl-proxy/pkg/log"
	"shorturl-proxy/services"
	"shorturl-proxy/services/shorturl"
	"shorturl-proxy/services/shorturl/proto"
)

type proxy struct {
	config *config.Config
	log    log.ILogger
}

func NewProxy(config *config.Config, log log.ILogger) *proxy {
	return &proxy{
		config: config,
		log:    log,
	}
}
func (p *proxy) PublicProxy(ctx *gin.Context) {
	p.redirection(ctx, true)
}
func (p *proxy) UserProxy(ctx *gin.Context) {
	p.redirection(ctx, false)
}

func (p *proxy) redirection(ctx *gin.Context, isPublic bool) {
	shortKey := ctx.Param("short_key")
	originalUrl, err := p.getOriginalUrl(shortKey, isPublic)
	if err != nil {
		p.log.Error(err)
		ctx.JSON(http.StatusInternalServerError, nil)
		return
	}
	ctx.Redirect(http.StatusFound, originalUrl)
}

func (p *proxy) getOriginalUrl(shortKey string, isPublic bool) (string, error) {
	pool := shorturl.NewShortUrlClientPool()
	conn := pool.Get()
	defer pool.Put(conn)
	client := proto.NewShortUrlClient(conn)
	ctx := services.AppendBearerToken(context.Background(), p.config.DependOn.ShortUrl.AccessToken)
	res, err := client.GetOriginUrl(ctx, &proto.ShortKey{Key: shortKey, UserID: 0, IsPublic: isPublic})
	if err != nil {
		p.log.Error(err)
		return "", err
	}
	return res.Url, err

}
