package main

const NODES_TOTAL int = 5

func main() {

	for i := range NODES_TOTAL {
		n := Node{id: i}

		go n.init(NODES_TOTAL)
	}

	select {}
}
