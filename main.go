package main

import (
	"github.com/magicspace/supernode/p2p"
	"github.com/magicspace/supernode/utils"
)

func main() {

	//ctx := context.Background()

	_, err := p2p.MakeNode()

	utils.HandleError(err, "Failed to initialize node", true)

	//go p2p.PeerDiscovery(ctx, host)

	select {}
}
