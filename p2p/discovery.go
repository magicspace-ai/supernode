package p2p

import (
	"context"
	"fmt"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/magicspace/supernode/utils"
)

// Discover peers using the boostrap method
// This helps other nodes on the network find themseleves
func DiscoverPeers(
	ctx context.Context, 
	dht *dht.IpfsDHT,
	rHost *rhost.RoutedHost,
) {

	fmt.Println("Announcing our node's presence")

	peerDiscoveryName := utils.GetConfig(
							"node.peerDiscoveryName", 
							"magicspace.ai",
						).(string)

	routingDiscovery := drouting.NewRoutingDiscovery(dht)
	dutil.Advertise(ctx, routingDiscovery, peerDiscoveryName)

	anyConnected := false
	for !anyConnected {
		
		fmt.Println("Searching for peers...")

		peerChan, err := routingDiscovery.FindPeers(ctx, peerDiscoveryName)
		
		if err != nil {
			utils.HandleError(err, "", true)
		}

		for peerInfo := range peerChan {
			if peerInfo.ID == rHost.ID() {
				continue // No self connection
			}
			
			err := rHost.Connect(ctx, peerInfo)
			if err != nil {
				fmt.Println("Failed connecting to ", peerInfo.ID.Pretty(), ", error:", err)
			} else {
				fmt.Println("Connected to:", peerInfo.ID.Pretty())
				anyConnected = true
			}
		}
	}

	fmt.Println("Peer discovery complete")
}	

