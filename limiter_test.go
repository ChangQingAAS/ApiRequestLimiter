package main

import (
	"fmt"
	"github.com/ChangQingAAS/ApiRequestLimiter/file"
	"github.com/ChangQingAAS/ApiRequestLimiter/limiter"
	"github.com/ChangQingAAS/ApiRequestLimiter/random"
	"runtime/debug"
	"strconv"
	"sync"
	"testing"
	"time"
)

func init() {
	filename := "./logger/log.log" // 日志文件路径
	file.Del(filename)
}

func TestLimiterAgent_HandleRequest(t *testing.T) {
	userNumber := 100
	requestNumber := 100
	testN(t, userNumber, requestNumber)
}

func testN(t *testing.T, userNumber int, requestNumber int) {
	w := GoN(userNumber, func(i int, user string) {
		for j := 0; j < requestNumber; j++ {
			time.Sleep(1000 * time.Millisecond)
			result, err := limiter.GLimiterAgent().HandleRequest(user, int64(random.RandInt(0, requestNumber)))
			if err != nil {
				t.Log(i, result, err)
			}
		}
	})
	w()
}

// GoN 同时启动多个协程，返回等待函数
func GoN(n int, fn func(int, string)) func() {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			user := "user" + strconv.Itoa(i)
			time.Sleep(1000 * time.Millisecond)
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(fmt.Sprintf("panic %s\n", err))
					fmt.Println(fmt.Sprint(string(debug.Stack())))
				}
			}()
			fn(i, user)
			wg.Done()
		}(i)
	}
	return wg.Wait
}

func BenchmarkLimiterAgent_HandleRequest(b *testing.B) {
	requestNumber := 100
	re := make(chan int)
	for i := 0; i < b.N; i++ {
		user := "user" + strconv.Itoa(i)
		time.Sleep(1000 * time.Millisecond)
		//go request(random.RandInt(0, userNumber), t, re, user, requestNumber)
		go request(i, b, re, user, requestNumber)
	}
	for i := 0; i < b.N*requestNumber; i++ {
		<-re
		//data := <-re
		//fmt.Println("i from re is ", data)
	}
}

func request(i int, b *testing.B, re chan int, user string, requestNumber int) {
	for j := 0; j < requestNumber; j++ {
		time.Sleep(1000 * time.Millisecond)
		result, err := limiter.GLimiterAgent().HandleRequest(user, int64(j))
		//re <- result
		re <- i*requestNumber + j
		if err != nil {
			b.Log(i, result, err)
		}
		//fmt.Println("put re ", i*requestNumber+j)
	}
}
