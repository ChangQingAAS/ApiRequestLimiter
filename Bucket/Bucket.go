package Bucket

import (
	"sync/atomic"
	"time"
)

// Bucket 令牌桶配置
type Bucket struct {
	Max   int64 //令牌桶的最大存储上限
	Cycle int64 //生成一批令牌的周期（每{cycle}毫秒生产一批令牌）
	Batch int64 //每批令牌的数量

	Residue int64 //令牌桶剩余空间
}

// NewTokenLimiter 初始化令牌桶全局限流器
func (bucket *Bucket) NewTokenLimiter() {
	//初始化令牌桶的剩余空间
	bucket.Residue = bucket.Max
	go func() {
		for {
			// 间隔一段时间发放令牌
			time.Sleep(time.Millisecond * time.Duration(bucket.Cycle))
			// 如果令牌数未超过上限，则继续累加
			if bucket.Residue+bucket.Batch <= bucket.Max {
				atomic.AddInt64(&bucket.Residue, bucket.Batch)
				continue
			} else {
				// 如果令牌数超过上限，则将令牌数设为上限值max
				atomic.StoreInt64(&bucket.Residue, bucket.Max)
			}
		}
	}()
}

// GetToken 获取令牌
// @num:本次请求需要拿取的令牌数
func (bucket *Bucket) GetToken(num int64) bool {

	//如果令牌桶剩余令牌数量不够
	if bucket.Residue-num <= 0 {
		return false
	}
	//令牌数量充足，取出令牌
	atomic.AddInt64(&bucket.Residue, -num)
	return true
}
