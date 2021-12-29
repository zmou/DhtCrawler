package dht

import (
	"DhtCrawler/krpc"
	"fmt"
	"io"
	"log"
)

type DhtNode struct {
	Node    *krpc.KNode
	Table   *krpc.KTable
	Network *Network
	Log     *log.Logger
	Master  chan string
	Krpc    *krpc.KRPC
	OutChan chan string
}

func NewDhtNode(id *Id, logger io.Writer, outHashIdChan chan string, master chan string) *DhtNode {
	node := new(krpc.KNode)
	node.Id = *id

	dht := new(DhtNode)
	dht.OutChan = outHashIdChan
	dht.Log = log.New(logger, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	dht.Node = node
	dht.Table = new(krpc.KTable)
	dht.Network = NewNetwork(dht)
	dht.Krpc = krpc.NewKRPC(dht)
	dht.Master = master

	return dht
}

func (dhtNode *DhtNode) Run() {

	//当前DHT节点运转进程
	go func() { dhtNode.Network.Listening() }()

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
