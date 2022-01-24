package dht

import (
	"fmt"
	"io"
	"log"
)

type DhtNode struct {
	Node    *KNode
	Table   *KTable
	Network *Network
	Log     *log.Logger
	Master  chan string
	Krpc    *KRPC
	OutChan chan string
	Port    int
}

func NewDhtNode(id *Id, logger io.Writer, outHashIdChan chan string, master chan string, port int) *DhtNode {
	node := new(KNode)
	node.Id = *id

	dht := new(DhtNode)
	dht.OutChan = outHashIdChan
	dht.Log = log.New(logger, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	dht.Node = node
	dht.Port = port
	dht.Table = new(KTable)
	dht.Network = NewNetwork(dht)
	dht.Krpc = NewKRPC(dht)
	dht.Master = master

	return dht
}

func (dhtNode *DhtNode) Run() {

	//当前DHT节点运转进程
	go func() { dhtNode.Listening() }()

	//自动结交更多DHT节点进程进程
	go func() { dhtNode.NodeFinder() }()

	dhtNode.Log.Println(fmt.Sprintf("DhtCrawler %s is runing...", dhtNode.Network.Conn.LocalAddr().String()))

	for {
		select {
		case msg := <-dhtNode.Master:
			dhtNode.Log.Println(msg)
		}
	}
}
