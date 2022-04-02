### 记录一下自己踩坑的经历

1. 因为之前工作上有使用copydata的需求所以找到了这个库，但是使用的时候发现接收数据时会乱码，所以做了编码的处理。之后又发现之前的版本会阻塞协成没有释放方法不符合我的使用需求，然后补上了退出的方法：
```go
func (p *BackWnd) SendThreadCloseMessage(idThread uint32) {
	win.PostThreadMessage(idThread, win.WM_QUIT, 0, 0)
}
```
只要在获取到关闭信号的时候执行 **goback.SendThreadCloseMessage(idThread)** 就可以退出了

顺便还加了另一种关闭的方法：
```go
func (p *BackWnd) SyncWaitMessage(quitCh chan int) {
    var msg win.MSG
    for {
        select {
        case <-quitCh:
            win.PostQuitMessage(0)
            log.Println("win.PostQuitMessage send")
        default:
            _ = win.PeekMessage(&msg, 0, 0, 0, 1)
            win.TranslateMessage(&msg)
            win.DispatchMessage(&msg)
            if msg.Message == win.WM_QUIT {
                log.Println("accept wm_quit message")
                return
            }
        }
    }
}
```
因为感兴趣就多研究了一下，其实两种功能上没啥区别，一定要说区别的话就是前者是在其他协成发出的关闭消息，而后者是在同一个协成发送的关闭消息了。具体的可以看看这个[用例](https://github.com/fanyiguang/Copydata/blob/master/go/src/example/main.go)

其他的就不多说了下面两个仓库里面都有：
1. [chenfeng8742](https://github.com/chenfeng8742/goback)
2. [jacky2478](https://github.com/jacky2478/goback)


