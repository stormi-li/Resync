package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	resync "github.com/stormi-li/Resync"
)

var redis_addr = "your redis addr"

func main() {
	go func() {
		LockProcess("1")
	}()
	go func() {
		LockProcess("2")
	}()
	go func() {
		LockProcess("3")
	}()
	time.Sleep(4 * time.Second)
}

func LockProcess(msg string) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redis_addr,
	})
	client := resync.NewClient(redisClient)
	lock := client.NewLock("lock1")
	lock.Lock()
	fmt.Println(lock.IsValid())
	fmt.Println(msg)
	time.Sleep(1 * time.Second)
	lock.Unlock()
}
