package main

import (
	"fmt"
	"time"

	resync "github.com/stormi-li/Resync"
)

func main() {
	client, _ := resync.NewClient("118.25.196.166:6379")
	client.SetNameSpace("fsdf")
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
