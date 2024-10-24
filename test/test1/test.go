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
	client.SetNamespace("fsdf")
	go func() {
		lock := client.NewLock("lock1")
		lock.Lock()
		fmt.Println(lock.IsValid())
		fmt.Println("1")
		time.Sleep(1 * time.Second)
		lock.Unlock()
	}()
	go func() {
		lock := client.NewLock("lock1")
		lock.Lock()
		fmt.Println(lock.IsValid())
		fmt.Println("2")
		time.Sleep(1 * time.Second)
		lock.Unlock()
	}()
	go func() {
		lock := client.NewLock("lock1")
		lock.Lock()
		fmt.Println(lock.IsValid())
		fmt.Println("3")
		time.Sleep(1 * time.Second)
		lock.Unlock()
	}()
	select {}
}
