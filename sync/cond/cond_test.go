package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	signalCh = make(chan os.Signal, 1)
	wg       sync.WaitGroup
	status   uint8
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		<-signalCh
		cancel()

		fmt.Println("信号监听 goroutine 退出")
	}(&wg)

	cond := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(ctx context.Context, index int, wg *sync.WaitGroup) {
			defer wg.Done()

			ticker := time.NewTicker(1 * time.Second)

		out:
			for {
				select {
				case <-ctx.Done():
					break out
				case <-ticker.C:
					cond.L.Lock()
					for status == 0 {
						fmt.Printf("---------->[%d], 等待中...\n", index)
						cond.Wait()
					}
					cond.L.Unlock()

					if status == 2 {
						break out
					}

					fmt.Printf("==========>[%d], 执行................!!!\n", index)

				}
			}
			fmt.Printf("%d 号 goroutine 退出...\n", index)
		}(ctx, i, &wg)
	}

	var (
		startTimer = time.NewTimer(5 * time.Second)
		pauseTimer = time.NewTimer(10 * time.Second)
	)

out:
	for {
		select {
		case <-ctx.Done():
			cond.L.Lock()
			status = 2
			cond.L.Unlock()
			cond.Broadcast()
			break out
		case <-startTimer.C:
			fmt.Println("发送启动信号...")
			pauseTimer.Reset(5 * time.Second)
			cond.L.Lock()
			status = 1
			cond.L.Unlock()
			cond.Broadcast()
		case <-pauseTimer.C:
			fmt.Println("发送暂停信号...")
			startTimer.Reset(5 * time.Second)
			cond.L.Lock()
			status = 0
			cond.L.Unlock()
			cond.Broadcast()
		}
	}

	fmt.Println("主进程等待所有 goroutine 退出...")
	wg.Wait()
}
