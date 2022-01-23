package main

import (
	"DhtCrawler/dht"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	stopChan = make(chan struct{}, 1)
	//爬虫输出抓取到的hashIds通道
	outHashIdChan = make(chan string, 100)
)

func main() {
	defer func() {
		<-stopChan
	}()

	//主进程
	master := make(chan string)

	//开启的dht节点
	for i := 0; i < 2; i++ {
		go func() {
			id := dht.GenerateID()
			dhtNode := dht.NewDhtNode(&id, os.Stdout, outHashIdChan, master)

			dhtNode.Run()
		}()
	}

	go func() {
		for {
			infoHash, ok := <-outHashIdChan
			if !ok {
				break
			}

			fmt.Println(infoHash)

			// 写入文件
			writeToFile(infoHash)
		}
	}()

	go func() {
		for {
			msg := <-master
			fmt.Println(msg)
		}
	}()
}

func writeToFile(hashId string) {
	logPath := "info_hash/" + time.Now().Format("2006-01-02") + ".log"
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("文件写入错误", err)
		return
	}

	hashId = fmt.Sprintf("%s\r\n", hashId)

	_, err = io.WriteString(f, hashId)

	defer func(f *os.File) {
		_ = f.Close()
	}(f)
}
