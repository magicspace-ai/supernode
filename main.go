package main

import (
	"context"

	"github.com/magicspace/supernode/p2p"
	"github.com/magicspace/supernode/utils"
)


func main() {

	ctx := context.Background()

	// initialize the node
	rhost, kDHT, err := p2p.MakeNode(ctx)

	utils.HandleError(err, "Node initialization failed", true)

	go p2p.DiscoverPeers(ctx, rhost, kDHT)

	select {}
}
