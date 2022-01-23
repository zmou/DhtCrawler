package dht

import (
	"fmt"
	"net"
	"time"

	"github.com/zeebo/bencode"
)

var BOOTSTRAP = []string{
	"router.utorrent.com:6881",
	"router.bittorrent.com:6881",
	"dht.transmissionbt.com:6881",
}

func (dhtNode *DhtNode) FindNode(node *KNode) {
	var id Id
	if node.Id != nil {
		id = node.Id.Neighbor()
	} else {
		id = dhtNode.Node.Id.Neighbor()
	}
	tid := dhtNode.Krpc.GenTID()
	v := make(map[string]interface{})
	v["t"] = fmt.Sprintf("%d", tid)
	v["y"] = "q"
	v["q"] = "find_node"
	args := make(map[string]string)
	args["id"] = string(id)
	args["target"] = string(GenerateID())
	v["a"] = args
	data, err := bencode.EncodeString(v)
	if err != nil {
		dhtNode.Log.Fatalln(err)
	}

	raddr := new(net.UDPAddr)
	raddr.IP = node.Ip
	raddr.Port = node.Port

	err = dhtNode.Network.Send([]byte(data), raddr)
	if err != nil {
		dhtNode.Log.Println(err)
	}
}

func (dhtNode *DhtNode) NodeFinder() {

	for {
		//	dhtNode.log.Println(len(dhtNode.table.Nodes), "port: ==== ", dhtNode.node.Port)

		if len(dhtNode.Table.Nodes) == 0 {
			for _, host := range BOOTSTRAP {
				raddr, err := net.ResolveUDPAddr("udp", host)
				if err != nil {
					dhtNode.Log.Fatalf("Resolve DNS error, %s\n", err)
					return
				}
				node := new(KNode)
				node.Port = raddr.Port
				node.Ip = raddr.IP
				node.Id = nil

				dhtNode.FindNode(node)
			}
		} else {
			for _, node := range dhtNode.Table.Nodes {
				dhtNode.FindNode(node)
			}
			dhtNode.Table.Nodes = nil
			time.Sleep(1 * time.Second)
		}

	}
}
