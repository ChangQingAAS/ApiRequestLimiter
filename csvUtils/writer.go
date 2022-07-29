package csvUtils

import (
	"encoding/csv"
	"log"
	"os"
)

// WriterCSV csv文件写入
// 依次为：请求ID,请求数量，新增数量,现有数量
func WriterCSV(path string, str []string) {

	//OpenFile读取文件，不存在时则创建，使用追加模式
	File, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("文件打开失败！")
	}
	defer File.Close()

	//创建写入接口
	WriterCsv := csv.NewWriter(File)
	//str := []string{"chen1", "hai1", "wei1"} //需要写入csv的数据，切片类型

	//写入一条数据，传入数据为切片(追加模式)
	err1 := WriterCsv.Write(str)
	if err1 != nil {
		log.Println("WriterCsv写入文件失败")
	}
	WriterCsv.Flush() //刷新，不刷新是无法写入的
	//log.Println("数据写入成功...")
}
