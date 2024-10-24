package main

import (
	"fmt"
	"time"

	resync "github.com/stormi-li/Resync"
)

func main() {
	client, _ := resync.NewClient("118.25.196.166:6379")
	client.SetNameSpace("fsdfs")
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
