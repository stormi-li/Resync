package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	resync "github.com/stormi-li/Resync"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "118.25.196.166:6379",
	})
	client := resync.NewClient(redisClient)
	go func() {
		lock := client.NewLock("lock1")
		lock.Lock()
		fmt.Println(lock.IsValid())
		fmt.Println("1111")
		time.Sleep(2 * time.Second)
		fmt.Println(lock.IsValid())
		lock.Unlock()
	}()
	time.Sleep(100 * time.Millisecond)
	select {}

}
