package logger

import (
	"io"
	"log"
	"os"
)

func WriterLog(str string) {
	// 创建、追加、读写，777，所有权限
	f, err := os.OpenFile("./logger/log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
	}()

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	log.Printf(str)
}
