package file

import (
	"fmt"
	"os"
)

func Del(file string) {
	err := os.Remove(file) //删除文件
	if err != nil {
		//如果删除失败则输出 file remove Error!
		fmt.Println("file remove Error!")
		//输出错误详细信息
		fmt.Printf("%s\n", err)
	} else {
		//如果删除成功则输出 file remove OK!
		fmt.Println("file remove OK!")
	}
}
