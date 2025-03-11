package shorturl

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mediahub/pkg/config"
	"mediahub/pkg/grpc_client_pool"
	"mediahub/pkg/log"
	"mediahub/pkg/zerror"
	"sync"
)

var pool grpc_client_pool.ClientPool
var once sync.Once

func NewShortUrlClientPool() grpc_client_pool.ClientPool {
	var err error
	if pool != nil {
		return pool
	}
	once.Do(func() {
		cnf := config.GetConfig()
		pool, err = grpc_client_pool.NewClientPool(cnf.DependOn.ShortUrl.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Error(zerror.NewByErr(err))
		}
	})
	return pool
}
