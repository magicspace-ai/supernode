package main

import (
	"context"
	"fmt"

	"github.com/magicspace/supernode/p2p"
	"github.com/magicspace/supernode/utils"
)


func main() {

	ctx := context.Background()

	// initialize the node
	rhost, kDHT, err := p2p.MakeNode(ctx)

	utils.HandleError(err, "Node initialization failed", true)

	go p2p.DiscoverPeers(ctx, rhost, kDHT)

	//lets init pubsub
	_, globalTopic, err := p2p.InitPubSub(ctx, rhost)

	utils.HandleError(err, "pubsub init error", true)

	//time.Sleep(5 * time.Second)

	utils.PrintInfo("Publishing Global Topic Msg")
	 err = globalTopic.Publish(ctx, []byte(fmt.Sprintf("Hello I am %s", rhost.ID().Pretty())))

	 utils.HandleError(err, "pubsub publish error", false)

	select {}
}
