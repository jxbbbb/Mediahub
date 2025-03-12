package grpc_client_pool

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"shorturl-proxy/pkg/log"
	"shorturl-proxy/pkg/zerror"
	"sync"
)

type ClientPool interface {
	Get() *grpc.ClientConn
	Put(*grpc.ClientConn)
}

type clientPool struct {
	pool sync.Pool
}

func NewClientPool(target string, opts ...grpc.DialOption) (ClientPool, error) {
	return &clientPool{
		pool: sync.Pool{
			New: func() any {
				conn, err := grpc.Dial(target, opts...) //...将切片打散
				if err != nil {
					log.Error(zerror.NewByErr(err))
					return nil
				}
				return conn
			},
		},
	}, nil
}
func (clientPool *clientPool) Get() *grpc.ClientConn {
	conn := clientPool.pool.Get().(*grpc.ClientConn)
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		conn = clientPool.pool.New().(*grpc.ClientConn)
	}
	return conn
}

func (clientPool *clientPool) Put(conn *grpc.ClientConn) {
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		return
	}
	clientPool.pool.Put(conn)
}
