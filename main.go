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

	dhtNodePorts = []int{
		38663,
		38723,
		38726,
		33968,
		55721,
		51886,
		56932,
		56324,
		52413,
		53938,
		57739,
		36949,
		37495,
		34605,
		48005,
	}
)

func main() {
	defer func() {
		<-stopChan
	}()

	//主进程
	master := make(chan string)

	//开启的dht节点
	for i := 0; i < 2; i++ {
		go func(i int) {
			id := dht.GenerateID()
			dhtNode := dht.NewDhtNode(&id, os.Stdout, outHashIdChan, master, dhtNodePorts[i])

			dhtNode.Run()
		}(i)
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
	rootPath, _ := os.Getwd()
	logPath := rootPath + "/info_hash/" + time.Now().Format("2006-01-02") + ".log"

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
