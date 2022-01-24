package dht

import (
	"net"
)

type Network struct {
	Conn *net.UDPConn
}

func NewNetwork(dhtNode *DhtNode) *Network {
	nw := new(Network)

	nw.Init(dhtNode)

	return nw
}

func (nw *Network) Init(dhtNode *DhtNode) {
	addr := new(net.UDPAddr)
	addr.Port = dhtNode.Port

	var err error
	nw.Conn, err = net.ListenUDP("udp", addr)

	if err != nil {
		panic(err)
	}

	laddr := nw.Conn.LocalAddr().(*net.UDPAddr)
	dhtNode.Node.Ip = laddr.IP
	dhtNode.Node.Port = laddr.Port
}

func (dhtNode *DhtNode) Listening() {
	buf := make([]byte, 1000)
	for {
		_, raddr, err := dhtNode.Network.Conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		_ = dhtNode.Decode(string(buf), raddr)
	}
}

func (dhtNode *DhtNode) Send(m []byte, addr *net.UDPAddr) error {
	_, err := dhtNode.Network.Conn.WriteToUDP(m, addr)

	if err != nil {
		dhtNode.Log.Println(err)
	}

	return err
}
