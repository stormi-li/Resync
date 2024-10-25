package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	resync "github.com/stormi-li/Resync"
)

var redisAddr = "118.25.196.166:3934"
var password = "12982397StrongPassw0rd"

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	client := resync.NewClient(redisClient, "lock-namespace")
	l := client.NewLock("mylock")
	l.Lock()
	fmt.Println("process2 locked")
	time.Sleep(5 * time.Second)
	l.Unlock()
	fmt.Println("process2 unlocked")
}
