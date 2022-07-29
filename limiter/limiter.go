package limiter

import (
	"fmt"
	"github.com/ChangQingAAS/ApiRequestLimiter/logger"
	"math"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

// HandleRequest 处理请求: 先获得锁，再处理
func (l *LimiterAgent) HandleRequest(user string, numRequest int64) (bool, error) {
	lockKey := getLimiterLockKey(user)

	// 获得锁
	for {
		lock, err := l.limiterGetLock(lockKey)
		if err != nil {
			return false, err
		}
		if lock {
			break
		}
	}

	finished, err := l.DoLimit(user, numRequest, time.Now().UnixNano())
	if finished != true {
		return false, err
	}

	// 解锁
	if err := l.limiterUnLock(lockKey); err != nil {
		return false, err
	}

	return true, nil
}

func (l *LimiterAgent) DoLimit(user string, numRequest int64, currNanoSec int64) (bool, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	conn := l.pool.Get()
	//fmt.Println("Do conn: ", conn)
	defer conn.Close()

	key := getLimiterKey(user)
	maxPermits, _ := strconv.ParseInt(l.spec.MaxPermits, 10, 64)
	rate, _ := strconv.ParseInt(l.spec.Rate, 10, 64)
	finished, err := l.Do(conn, user, key, currNanoSec, numRequest, maxPermits, rate)

	if finished != true {
		if err != nil {
			str := fmt.Sprintf("Can't handle: %s's %d requests, error: %v\n", user, numRequest, err)
			str += fmt.Sprintf("\t\t\t\t\t\t\tDetails: \n")
			str += fmt.Sprintf("\t\t\t\t\t\t\t\tCurrent Requests: %d\n", numRequest)
			str += fmt.Sprintf("\t\t\t\t\t\t\t\tUser Info: username: %s, Max: %s, Rate: %s\n", user, l.spec.MaxPermits, l.spec.Rate)
			str += fmt.Sprintf("\t\t\t\t\t\t\t\tBucket Info: Residue: %d, Max: %d, Cycle: %d, Batch: %d\n", centerBucket.Max, centerBucket.Max, centerBucket.Cycle, centerBucket.Batch)
			logger.WriterLog(str)
		}
		return false, err
	}
	return true, nil
}

func (l *LimiterAgent) Do(conn redis.Conn, user string, key string, currNanoSec int64, numRequest int64, maxPermits int64, rate int64) (finished bool, err error) {
	// 取当前纳秒数和token数量，若没有则初始化
	lastNanoSec, _ := redis.String(conn.Do("HGET", key, "lastNanoSec"))
	currPermits, _ := redis.String(conn.Do("HGET", key, "currPermits"))
	if lastNanoSec == "" {
		conn.Do("HSET", key, "lastNanoSec", currNanoSec)
		conn.Do("HSET", key, "currPermits", maxPermits-numRequest)
		currPermits = strconv.FormatInt(maxPermits-numRequest, 10)
		lastNanoSec = strconv.FormatInt(currNanoSec, 10)
	}

	// 计算当前保留了多少个token
	lastNanoSecInt64, _ := strconv.ParseInt(lastNanoSec, 10, 64)
	reservePermits := math.Ceil(float64((currNanoSec - lastNanoSecInt64) / int64(math.Pow(10, 9)) * rate))

	// 保留上次访问时间
	conn.Do("HSET", key, "lastNanoSec", currNanoSec)

	// 计算当前数量
	var current float64
	currPermitsFloat64, _ := strconv.ParseFloat(currPermits, 64)
	if reservePermits+currPermitsFloat64 > float64(maxPermits) {
		reservePermits = float64(maxPermits) - currPermitsFloat64
		current = float64(maxPermits)
	} else {
		current = reservePermits + currPermitsFloat64
	}

	// 处理请求，并写入处理请求后的token数量
	remaining := current - float64(numRequest)
	if remaining >= 0 {
		_, err = conn.Do("HSET", key, "currPermits", remaining)
		str := fmt.Sprintf("User Handle: %s successly handle %d requests\n", user, numRequest)
		str += fmt.Sprintf("\t\t\t\t\t\t\tDetails: \n")
		str += fmt.Sprintf("\t\t\t\t\t\t\t\tCurrent Requests: %d\n", numRequest)
		str += fmt.Sprintf("\t\t\t\t\t\t\t\tUser Info: username: %s, resvervePermits: %f, currPermits: %f, Max: %s, Rate: %s\n", user, reservePermits, current, l.spec.MaxPermits, l.spec.Rate)
		str += fmt.Sprintf("\t\t\t\t\t\t\t\tBucket Info: Residue: %d, Max: %d, Cycle: %d, Batch: %d\n", centerBucket.Max, centerBucket.Max, centerBucket.Cycle, centerBucket.Batch)
		logger.WriterLog(str)
		if err != nil {
			return false, err
		}
		return true, nil
	} else {
		// 本地用户的token不够了，向全局请求
		isOk := centerBucket.GetToken(numRequest)
		// 全局请求处理失败
		if isOk == false {
			// 这个请求处理不了
			str := fmt.Sprintf("Can't handle: %s's %d requests, because token in %s's buecket or token in centerBucket isn't enough!\n", user, numRequest, user)
			str += fmt.Sprintf("\t\t\t\t\t\t\tDetails: \n")
			str += fmt.Sprintf("\t\t\t\t\t\t\t\tCurrent Requests: %d\n", numRequest)
			str += fmt.Sprintf("\t\t\t\t\t\t\t\tUser Info: username: %s, resvervePermits: %f, currPermits: %f, Max: %s, Rate: %s\n", user, reservePermits, current, l.spec.MaxPermits, l.spec.Rate)
			str += fmt.Sprintf("\t\t\t\t\t\t\t\tBucket Info: Residue: %d, Max: %d, Cycle: %d, Batch: %d\n", centerBucket.Max, centerBucket.Max, centerBucket.Cycle, centerBucket.Batch)
			logger.WriterLog(str)
			_, err = conn.Do("HSET", key, "currPermits", current)
			if err != nil {
				return false, err
			}
		} else {
			// 向全局请求成功
			str := fmt.Sprintf("center Handle: %s get help from centerBucket to handle %d resuests!\n", user, numRequest)
			str += fmt.Sprintf("\t\t\t\t\t\t\tDetails: \n")
			str += fmt.Sprintf("\t\t\t\t\t\t\t\tCurrent Requests: %d\n", numRequest)
			str += fmt.Sprintf("\t\t\t\t\t\t\t\tUser Info: username: %s, resvervePermits: %f, currPermits: %f, Max: %s, Rate: %s\n", user, reservePermits, current, l.spec.MaxPermits, l.spec.Rate)
			str += fmt.Sprintf("\t\t\t\t\t\t\t\tBucket Info: Residue: %d, Max: %d, Cycle: %d, Batch: %d\n", centerBucket.Max, centerBucket.Max, centerBucket.Cycle, centerBucket.Batch)
			logger.WriterLog(str)
			_, err = conn.Do("HSET", key, "currPermits", current)
			if err != nil {
				return false, err
			}
			return true, nil
		}
		return false, nil
	}
}

// 上锁
func (l *LimiterAgent) limiterGetLock(lockKey string) (bool, error) {
	conn := l.pool.Get()
	//fmt.Println("Lock conn: ", conn)
	_, err := redis.String(conn.Do("SET", lockKey, 1, "EX", 1, "NX"))
	if err == redis.ErrNil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// 解锁
func (l *LimiterAgent) limiterUnLock(lockKey string) error {
	conn := l.pool.Get()
	//fmt.Println("Unlock conn: ", conn)
	_, err := conn.Do("DEL", lockKey)
	if err != nil {
		return err
	}
	return nil
}
