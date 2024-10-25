package resync

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	ripc "github.com/stormi-li/Ripc"
)

type Client struct {
	redisClient *redis.Client
	ripcClient  *ripc.Client
	namespace   string
	ctx         context.Context
}

const lockPrefix = "stormi:lock:"

func NewClient(redisClient *redis.Client, namespace string) *Client {
	ripcClient := ripc.NewClient(redisClient, namespace)
	return &Client{
		redisClient: redisClient,
		ripcClient:  ripcClient,
		namespace:   namespace + ":" + lockPrefix,
		ctx:         context.Background(),
	}
}

func (c *Client) NewLock(lockName string) *Lock {
	return &Lock{
		uuid:        uuid.NewString(),
		lockName:    lockName,
		stop:        make(chan struct{}, 1),
		ripcClient:  c.ripcClient,
		redisClient: c.redisClient,
		namespace:   c.namespace,
		ctx:         c.ctx,
	}
}
