package conf

import "strconv"

type Redis struct {
	Host string `json:"host"`
}

type Limiter struct {
	MaxPermits string `json:"maxPermits"`
	Rate       string `json:"rate"`
}

var (
	redis   Redis
	limiter Limiter
)

func init() {
	redis.Host = ":6379"
	limiter.MaxPermits = strconv.Itoa(500)
	limiter.Rate = strconv.Itoa(50)
}

func GetRedis() Redis {
	return redis
}

func GetLimiter() Limiter {
	return limiter
}
