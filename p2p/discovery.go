package p2p

import (
	"context"
	"fmt"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	network "github.com/libp2p/go-libp2p/core/network"
	peer "github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/magicspace/supernode/utils"
)

// Discover peers using the boostrap method
// This helps other nodes on the network find themseleves
func DiscoverPeers(
	ctx context.Context,
	rh  *rhost.RoutedHost,
	dht *dht.IpfsDHT,
) {

	fmt.Println("Announcing our node's presence")

	peerDiscoveryName := utils.GetConfig(
							"node.peerDiscoveryName", 
							"magicspace.ai",
						).(string)

	routingDiscovery := drouting.NewRoutingDiscovery(dht)
	dutil.Advertise(ctx, routingDiscovery, peerDiscoveryName)

	var peerConnectFailCount = make(map[peer.ID]int64)
	maxPeerConnectRetries := utils.GetConfig("node.maxPeerConnectRetries", 10).(int64)

	for {
		
		fmt.Println("Searching for peers...")

		peerChan, err := routingDiscovery.FindPeers(ctx, peerDiscoveryName)
		
		if err != nil {
			utils.HandleError(err, "", true)
		}

		for peerinfo := range peerChan {
			if peerinfo.ID == rh.ID() {
				continue // No self connection
			}	

			// if already connected, skip it
			if rh.Network().Connectedness(peerinfo.ID) == network.Connected {
				continue
			}
			
			if peerConnectFailCount[peerinfo.ID] > maxPeerConnectRetries {
				continue		
			}

			err := rh.Connect(ctx, peerinfo)
			//_, err := rh.Network().DialPeer(ctx, peerInfo.ID)

			if err != nil {
				fmt.Println("peer connection failed: ", peerinfo.ID.Pretty(), ", error:", err)
				peerConnectFailCount[peerinfo.ID] += 1
			} else {
				utils.PrintSuccess("Connected to: %s", peerinfo.ID.Pretty())
				peerConnectFailCount[peerinfo.ID] = 0
			}
		}

		time.Sleep(5 * time.Second)
	}

}	

