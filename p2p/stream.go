package p2p

import (
	"bufio"

	"github.com/libp2p/go-libp2p/core/network"
)

func streamHandler(stream network.Stream) {

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(rw)
	go writeData(rw)
}

/**
 * read incoming stream
 **/
 func readData(rw *bufio.ReadWriter) {
	//for {

//	}
}

/**
 * write incoming stream
 **/
func writeData(rw *bufio.ReadWriter) {
	//for {
		
	//}
}
