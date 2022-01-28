package main

import (
	"fmt"
	"goback/go/src/goback"
	"log"

	"golang.org/x/text/encoding/simplifiedchinese"
)

// MainWnd为UI主窗口名称
// 修复c++端访问goback只能在初始化时访问的bug，需要在隐藏窗口内部等待消息，
// 而不是后台服务example.exe循环等待从隐藏窗口推送过来的管道消息
func main() {
	obj := goback.Regist("MainWnd")
	go func() {
		for {
			msg, ok := <-obj.BufCh
			if !ok {
				break
			}
			fmt.Println("转码前结果：", msg)
			gbkMsg, err := simplifiedchinese.GB18030.NewDecoder().Bytes(msg)
			if err != nil {
				log.Println("decode failed:", err)
				return
			}

			fmt.Println("转码前结果:", string(gbkMsg))
		}
		close(obj.BufCh)
	}()
	goback.Wait()
}
