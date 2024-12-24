package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const PADDING = 2000

type Node struct {
	nodes map[int]int
	id    int

	arbitraryTime time.Duration
}

type Args struct {
	Clock   map[int]int
	NodeId  int
	Message string
}

type Response struct {
}

func RandomTimeout() time.Duration {
	return time.Duration((rand.Intn(300-150) + 150)) * time.Millisecond
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (n *Node) sendMessage(node_id int) {
	n.nodes[n.id]++

	address := fmt.Sprintf(":%d", PADDING+node_id)
	client, e := rpc.DialHTTP("tcp", address)

	if e != nil {
		log.Println("Error on dialing, ", e)
		return
	}

	args := Args{Clock: n.nodes, NodeId: n.id, Message: "Hey"}
	reply := Response{}

	log.Printf("[%d] sends Message to %d", n.id, node_id)
	client.Call("Node.ReceiveMessage", &args, &reply)
}

func (n *Node) ReceiveMessage(args *Args, response *Response) error {
	for k, v := range n.nodes {
		n.nodes[k] = Max(v, args.Clock[k])
	}

	n.nodes[n.id]++
	log.Printf("[%d] receives message. Clock: %v", n.id, n.nodes)
	return nil
}

func (n *Node) init(total int) {

	n.arbitraryTime = RandomTimeout()
	n.nodes = make(map[int]int)

	for i := range total {
		n.nodes[i] = 0
	}

	mux := http.NewServeMux()

	rpcServer := rpc.NewServer()
	err := rpcServer.Register(n)

	if err != nil {
		log.Println("Error on register, ", err)
		return
	}

	address := fmt.Sprintf(":%d", PADDING+n.id)

	listener, e := net.Listen("tcp", address)

	if e != nil {
		log.Println("error on listen ", e)
		return
	}

	mux.Handle("/_goRPC_", rpcServer)
	go http.Serve(listener, mux)

	for {
		time.Sleep(time.Second)

		if rand.Intn(10) > 7 {
			n.sendMessage(rand.Intn(total))
		}
	}

}
