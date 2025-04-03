package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	maxResponseSize = 10 * 1024 * 1024 // 限制响应大小为10MB
)

func sendRequest(ctx context.Context, url string, headers map[string]string, rest int, i int, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	defer func() { <-sem }()

	start := time.Now()

	// 创建一个可取消的请求上下文
	reqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 创建请求
	req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
	if err != nil {
		fmt.Printf("创建请求 %d 失败: %v\n\n", i+1, err)
		return
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 限制重定向的次数
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("停止过多重定向")
			}
			if len(via) > 0 {
				fmt.Printf("\n请求 %d 已重定向到: %s\n", i+1, req.URL.String())
			}
			return nil
		},
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("\n请求 %d 失败, 错误: %v\n", i+1, err)
		runtime.GC() // 手动触发垃圾回收
		time.Sleep(time.Duration(rest) * time.Second)
		return
	}
	defer resp.Body.Close()

	// 读取并丢弃响应体内容，但限制最大大小
	limitedReader := io.LimitReader(resp.Body, maxResponseSize)
	n, err := io.Copy(io.Discard, limitedReader)
	if err != nil {
		fmt.Printf("\n请求 %d 读取响应体失败: %v\n", i+1, err)
		return
	}

	// 打印成功请求的信息
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("\n\n请求 %d 成功, 最终地址: %s\n\n", i+1, resp.Request.URL.String())
		elapsed := time.Since(start).Seconds()

		var fileSize float64
		fileSizeStr := resp.Header.Get("Content-Length")
		if fileSizeStr != "" {
			sizeInt, err := strconv.Atoi(fileSizeStr)
			if err == nil {
				fileSize = float64(sizeInt) / 1024 / 1024
			} else {
				// 如果Content-Length解析失败，使用已读取的字节数
				fileSize = float64(n) / 1024 / 1024
			}
		} else {
			// 如果没有Content-Length header，使用已读取的字节数
			fileSize = float64(n) / 1024 / 1024
		}

		speed := fileSize / elapsed
		fmt.Printf("用时 %.2f 秒, 文件大小: %.2f MB, 速度: %.2f MB/s\n\n%s\n",
			elapsed, fileSize, speed, "----------------------------------------")
	} else {
		fmt.Printf("请求 %d 失败, 状态码: %d\n\n%s\n",
			i+1, resp.StatusCode, "----------------------------------------")
	}

	runtime.GC() // 手动触发垃圾回收
	time.Sleep(time.Duration(rest) * time.Second)
}

func main() {
	// 解析命令行参数
	times := flag.Int("t", 1, "请求次数")
	workers := flag.Int("w", 1, "并发线程数")
	ua := flag.String("ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0", "User-Agent")
	rest := flag.Int("r", 1, "请求间隔时间(秒)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("请提供请求地址")
		os.Exit(1)
	}
	url := args[0]

	// 设置请求头
	headers := map[string]string{
		"User-Agent": *ua,
	}

	// 限制最大并发数
	maxWorkers := *workers
	if maxWorkers <= 0 {
		maxWorkers = 1
	}

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxWorkers)

	// 发送请求
	fmt.Printf("开始执行 %d 次请求，最大并发数: %d\n", *times, maxWorkers)

	for i := 0; i < *times; i++ {
		sem <- struct{}{} // 获取信号量
		wg.Add(1)
		go sendRequest(ctx, url, headers, *rest, i, &wg, sem)
	}

	wg.Wait()

	// 清理资源
	close(sem)
	runtime.GC() // 手动触发垃圾回收

	fmt.Println("所有请求已完成，请检查结果。")
}
