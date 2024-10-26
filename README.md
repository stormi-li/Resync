# RESYNC Guides

Simple and secure distributed lock library.

# Overview

- Support watch dog
- Support lock id
- Support blocking and wake-up
- Every feature comes with tests
- Developer Friendly

# Install


```shell
go get -u github.com/stormi-li/Resync
```

# Quick Start

```go
package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	resync "github.com/stormi-li/Resync"
)

var redisAddr = “localhost:6379”
var password = “your password”

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	client := resync.NewClient(redisClient, "lock-namespace")
	l := client.NewLock("mylock")
	l.Lock()
	fmt.Println("process locked")
	time.Sleep(5 * time.Second)
	l.Unlock()
	fmt.Println("process unlocked")
}
```

# Interface - resync

## NewClient

### Create resync client
```go
package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	resync "github.com/stormi-li/Resync"
)

var redisAddr = “localhost:6379”
var password = “your password”

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
	})
	client := resync.NewClient(redisClient, "lock-namespace")
}
```

The first parameter is a redis client of successful connection, the second parameter is a unique namespace.

# Interface - resync.Client

## NewLock

### Create a distrubuted lock
```go
lock := client.NewLock("mylock")
```
The parameter is a unique lock name.

# Interface - resync.Lock

## Lock

### Preempt lock
```go
lock.Lock()
```
Failure to preempt a lock will block the process.

## Unlock 

### Release lock
```go
lock.Unlock()
```

# Community

## Ask

### How do I ask a good question?
- Email - 2785782829@qq.com
- Github Issues - https://github.com/stormi-li/Resync/issues