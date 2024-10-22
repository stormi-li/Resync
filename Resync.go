package resync

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	ripc "github.com/stormi-li/Ripc"
)

type Client struct {
	ripcClient *ripc.Client
}

func NewClient(addr string) (*Client, error) {
	client := Client{}
	ripcClient, err := ripc.NewClient(addr)
	if err != nil {
		return nil, err
	}
	client.ripcClient = ripcClient
	return &client, nil
}

type Lock struct {
	uuid        string
	lockName    string
	stop        chan struct{}
	ripcClient  *ripc.Client
	redisClient *redis.Client
	ctx         context.Context
}

func (client *Client) NewLock(lockName string) *Lock {
	lock := Lock{}
	lock.ripcClient = client.ripcClient
	lock.uuid = uuid.New().String()
	lock.lockName = lockName
	lock.stop = make(chan struct{}, 1)
	lock.ctx = context.Background()
	lock.redisClient = lock.ripcClient.RedisClient
	return &lock
}

func (l *Lock) Lock() {
	for {
		var ok bool
		ok, _ = l.redisClient.SetNX(l.ctx, l.lockName, l.uuid, 3*time.Second).Result()

		if ok {
			go func() {
				ticker := time.NewTicker(1 * time.Second)
				defer ticker.Stop()
				for {
					select {
					case <-ticker.C:
						l.updateExpiryIfValueMatches()
					case <-l.stop:
						return
					}
				}
			}()
			break
		} else {
			l.ripcClient.Wait(l.ctx, l.lockName, 3*time.Second)
		}
	}
}

func (l *Lock) Unlock() {
	l.stop <- struct{}{}
	l.deleteIfValueMatches()
	l.ripcClient.Notify(l.ctx, l.lockName, "unlock")

}

func (l *Lock) IsValid() bool {
	res, _ := l.updateExpiryIfValueMatches()
	return res
}

// updateExpiryIfValueMatches 更新过期时间，如果键的值匹配预期值
func (l *Lock) updateExpiryIfValueMatches() (bool, error) {
	script := `
        local key = KEYS[1]
        local expected_value = ARGV[1]
        local new_ttl = ARGV[2]

        local current_value = redis.call('GET', key)

        if current_value == expected_value then
            redis.call('EXPIRE', key, new_ttl)
            return 1
        else
            return 0
        end
    `

	result, err := l.redisClient.Eval(l.ctx, script, []string{l.lockName}, l.uuid, 3).Result()
	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}

// deleteIfValueMatches 判断键的值是否匹配目标值，若匹配则删除该键
func (l *Lock) deleteIfValueMatches() (bool, error) {
	luaScript := `
		local currentValue = redis.call("GET", KEYS[1])
		if currentValue == ARGV[1] then
			redis.call("DEL", KEYS[1])
			return 1  -- 1表示成功删除
		else
			return 0  -- 0表示值不匹配
		end
	`
	result, err := l.ripcClient.RedisClient.Eval(l.ctx, luaScript, []string{l.lockName}, l.uuid).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}

// type ReadWriteLock struct {
// 	writeLock *Lock
// 	readLock  *Lock
// }

// func (client *Client) newReadWriteLock(lockName string) *ReadWriteLock {
// 	lock := ReadWriteLock{}
// 	lock.writeLock = client.NewLock(lockName)
// 	lock.readLock = client.NewLock(lockName + "-count")
// 	return &lock
// }

// type WriteLock struct {
// 	writeLock *Lock
// 	readLock  *Lock
// }

// func (readWriteLock *ReadWriteLock) GetWriteLock() *WriteLock {
// 	lock := WriteLock{writeLock: readWriteLock.writeLock, readLock: readWriteLock.readLock}
// 	return &lock
// }

// func (l *WriteLock) Lock() {
// 	l.writeLock.Lock()
// 	l.readLock.Lock()
// }

// func (l *WriteLock) IsValid() bool {
// 	return l.writeLock.IsValid()
// }

// func (l *WriteLock) Unlock() {
// 	l.readLock.Unlock()
// 	l.writeLock.Unlock()
// }

// type ReadLock struct {
// 	writeLockName string
// 	readLockName  string
// 	stop          chan struct{}
// 	ripcClient    *ripc.Client
// 	redisClient   *redis.Client
// 	ctx           context.Context
// }

// func (l *ReadWriteLock) GeReadLock() *ReadLock {
// 	lock := ReadLock{
// 		writeLockName: l.writeLock.lockName,
// 		readLockName:  l.readLock.lockName,
// 		stop:          make(chan struct{}, 1),
// 		ripcClient:    l.readLock.ripcClient,
// 		redisClient:   l.readLock.redisClient,
// 		ctx:           l.readLock.ctx,
// 	}
// 	return &lock
// }

// func (l *ReadLock) Lock() {
// 	for {
// 		var ok bool
// 		ok, _ = l.IncrementOrSetIfNotExists()

// 		if ok {
// 			go func() {
// 				ticker := time.NewTicker(1 * time.Second)
// 				defer ticker.Stop()
// 				for {
// 					select {
// 					case <-ticker.C:
// 						l.redisClient.Expire(l.ctx, l.readLockName, 3*time.Second)
// 					case <-l.stop:
// 						return
// 					}
// 				}
// 			}()
// 			break
// 		} else {
// 			l.ripcClient.Wait(l.ctx, l.writeLockName, 3*time.Second)
// 		}
// 	}
// }

// func (l *ReadLock) IsValid() bool {
// 	str, _ := l.redisClient.Get(l.ctx, l.readLockName).Result()
// 	_, err := strconv.Atoi(str)
// 	return err == nil
// }

// func (l *ReadLock) Unlock() {
// 	l.stop <- struct{}{}
// 	count, _ := l.DecrementOrDelete()
// 	if count < 1 {
// 		l.ripcClient.Notify(l.ctx, l.readLockName, "unlock")
// 	}
// }

// func (l *ReadLock) DecrementOrDelete() (int, error) {
// 	script := `
//     local key = KEYS[1]
//     local currentValue = redis.call('GET', key)

//     if currentValue then
//         currentValue = tonumber(currentValue)
//         if currentValue > 1 then
//             redis.call('DECR', key)
//             return currentValue - 1
//         elseif currentValue == 1 then
//             redis.call('DEL', key)
//             return 0
//         end
//     end

//     return -1
//     `

// 	result, err := l.redisClient.Eval(l.ctx, script, []string{l.readLockName}).Result()
// 	return result.(int), err
// }

// func (l *ReadLock) IncrementOrSetIfNotExists() (bool, error) {
// 	luaScript := `
//     if redis.call("EXISTS", KEYS[1]) == 1 then
//         return false
//     else
//         local bValue = redis.call("INCR", KEYS[2])
//         if bValue == 1 then
//             redis.call("SET", KEYS[2], 1)
//         end
//         return true
//     end
//     `

// 	result, err := l.redisClient.Eval(l.ctx, luaScript, []string{l.writeLockName, l.readLockName}).Result()
// 	if err != nil {
// 		return false, err
// 	}

// 	return result.(int64) != 0, nil
// }
