package csvUtils

import (
	"encoding/csv"
	"log"
	"os"
)

func ReadCsv(filepath string) [][]string {
	//打开文件(只读模式)，创建io.read接口实例
	opencast, err := os.Open(filepath)
	if err != nil {
		log.Println("csv文件打开失败！")
	}
	defer opencast.Close()

	//创建csv读取接口实例
	ReadCsv := csv.NewReader(opencast)

	////获取一行内容，一般为第一行内容
	//read, _ := ReadCsv.Read() //返回切片类型：[chen  hai wei]
	//log.Println(read)

	//读取所有内容
	ReadAll, err := ReadCsv.ReadAll() //返回切片类型：[[s s ds] [a a a]]
	return ReadAll

	/*
	  说明：
	   1、读取csv文件返回的内容为切片类型，可以通过遍历的方式使用或Slicer[0]方式获取具体的值。
	   2、同一个函数或线程内，两次调用Read()方法时，第二次调用时得到的值为每二行数据，依此类推。
	   3、大文件时使用逐行读取，小文件直接读取所有然后遍历，两者应用场景不一样，需要注意。
	*/
}
