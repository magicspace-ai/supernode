package main

import (
	"context"

	"github.com/magicspace/supernode/p2p"
	"github.com/magicspace/supernode/utils"
)

func main() {

	ctx := context.Background()

	node, err := p2p.MakeNode()

	utils.HandleError(err, "Node initialization failed", true)

	go p2p.DiscoverPeers(ctx, node)

	select {}
}
