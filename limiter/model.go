package limiter

import (
	"errors"
	"github.com/ChangQingAAS/ApiRequestLimiter/Bucket"
	"github.com/ChangQingAAS/ApiRequestLimiter/conf"
	"github.com/ChangQingAAS/ApiRequestLimiter/data"
	"github.com/gomodule/redigo/redis"
	"sync"
)

var limiterAgent *LimiterAgent
var once sync.Once
var requestId int
var centerBucket Bucket.Bucket

type LimiterValue struct {
	MaxPermits string
	Rate       string
}

type LimiterAgent struct {
	spec LimiterValue
	pool *redis.Pool
	lock sync.Mutex
}

func init() {
	requestId = 0
	runCenterBucket()
}

func runCenterBucket() {
	//设置令牌桶最大容量为100000，每500毫秒生产5000个令牌，相当于每1秒最多只能取出10000个令牌
	centerBucket = Bucket.Bucket{
		Max:   100000,
		Cycle: 500,
		Batch: 5000,
	}

	//初始化令牌桶限流器
	centerBucket.NewTokenLimiter()
}

func GLimiterAgent() *LimiterAgent {
	once.Do(func() {
		limiterAgent = NewLimiterAgent()
	})
	return limiterAgent
}

func NewLimiterAgent() *LimiterAgent {
	if data.RedisPool() == nil {
		panic(errors.New("pool error"))
		return nil
	}

	limiterAgent := &LimiterAgent{
		spec: LimiterValue{
			MaxPermits: conf.GetLimiter().MaxPermits,
			Rate:       conf.GetLimiter().Rate,
		},
		pool: data.RedisPool(),
		lock: sync.Mutex{},
	}

	return limiterAgent
}
