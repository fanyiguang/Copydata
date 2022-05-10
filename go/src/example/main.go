package main

import (
	"context"
	"copydata/go/src/goback"
	"fmt"
	"log"
	"time"

	"golang.org/x/sys/windows"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func main() {
	cancel, cancelFunc := context.WithCancel(context.Background())
	go AsyncProcMode(cancel)
	time.Sleep(time.Second * 15)
	cancelFunc()
	log.Println("Async cpData end")

	cancel1, cancelFunc1 := context.WithCancel(context.Background())
	go SyncProcMode(cancel1)
	time.Sleep(time.Second * 15)
	cancelFunc1()
	log.Println("Sync cpData end")
	time.Sleep(time.Second * 5)
}

func AsyncProcMode(ctx context.Context) {
	var idThread uint32
	obj := goback.Regist("MainWnd")
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("cpData accept loop done")
				goback.SendThreadCloseMessage(idThread)
				return
			case msg, ok := <-obj.BufCh:
				if !ok {
					continue
				}
				fmt.Println("转码前结果：", msg)
				gbkMsg, err := simplifiedchinese.GB18030.NewDecoder().Bytes(msg)
				if err != nil {
					log.Println("decode failed:", err)
					return
				}

				fmt.Println("转码后结果:", string(gbkMsg))
			}
		}
	}()
	idThread = windows.GetCurrentThreadId()
	goback.Wait()
}

func SyncProcMode(ctx context.Context) {
	doneCh := make(chan int, 1)
	obj := goback.Regist("MainWnd")
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("cpData accept loop done")
				doneCh <- 1
				return
			case msg, ok := <-obj.BufCh:
				if !ok {
					continue
				}
				fmt.Println("转码前结果：", msg)
				gbkMsg, err := simplifiedchinese.GB18030.NewDecoder().Bytes(msg)
				if err != nil {
					log.Println("decode failed:", err)
					return
				}

				fmt.Println("转码后结果:", string(gbkMsg))
			}
		}
	}()
	goback.SyncWait(doneCh)
}
