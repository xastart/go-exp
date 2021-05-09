package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)
import "golang.org/x/sync/errgroup"

// HttpServerStart 启动http server
func HttpServerStart(ctx context.Context) error {
	server := &http.Server{Addr: ":8082", Handler: nil}

	go func() {
		select {
		case <-ctx.Done():
			server.Close()
			fmt.Println("cleanup httpServer1")
		}
	}()

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello!"))
	})

	fmt.Println("start httpServer1")
	return server.ListenAndServe()

}

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	// 启动httpserver1
	g.Go(func() error {
		return HttpServerStart(ctx)
	})

	// 模拟启动httpserver2
	g.Go(func() error {
		//
		fmt.Println("start server2")
		select {
		case <-ctx.Done():
			{
				fmt.Println("cleanup server2")
				// 模拟清理操作
				time.Sleep(2 * time.Second)
				return errors.New("serverExit")
			}
		}

	})

	// 注册系统事件
	g.Go(func() error {
		sysChan := make(chan os.Signal, 0)
		signal.Notify(sysChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		select {
		case s := <-sysChan:
			fmt.Println("sys signal:", s)
			//cancel() // 取消请求，通知用到ctx的所有goroutine
			return errors.New("killError")

		case <-ctx.Done():
			return errors.New("internalError")

		}

	})

	// 模拟其他异常错误，引发所有goroutine 全部退出
	//g.Go(func() error {
	//	return errors.New("testExit")
	//})

	// 等待goroutine完成
	if err := g.Wait(); err != nil {
		fmt.Printf("server error:%v", err)
	}

}
