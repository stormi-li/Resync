# Resync 框架

## 简介

Resync 是一个基于 Redis 的分布式锁框架，能够为分布式系统提供高效、安全的锁管理机制，确保多进程的并发操作安全。

## 功能

- 支持看门狗机制：持锁进程可以定期对锁续约，防止由于超时导致锁意外释放。
- 支持锁身份识别：进程可以在续约和释放锁时验证自己是否是锁的合法持有者。
- 支持阻塞和唤醒：抢占锁失败后能够进入阻塞状态，等待自动唤醒或被其它进程唤醒。

## 安装

```shell
go get github.com/stormi-li/Resync
```

## 使用

**示例代码**：

```go
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
	lock := client.NewLock("lock1") //创建一个叫lock1的分布式锁
	lock.Lock()
	fmt.Println(lock.IsValid()) //检查是否持有锁
	fmt.Println(msg)
	time.Sleep(1 * time.Second)
	lock.Unlock()
}
```



