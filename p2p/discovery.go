package p2p

import (
	"context"
	"fmt"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	network "github.com/libp2p/go-libp2p/core/network"
	peer "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/magicspace/supernode/utils"
)

// connected peers
var ConnectedPeers = make(map[peer.ID]peer.AddrInfo)
//var Swarm swarm.Swarm

// Discover peers using the boostrap method
// This helps other nodes on the network find themseleves
func DiscoverPeers(
	ctx context.Context,
	rh  *rhost.RoutedHost,
	dht *dht.IpfsDHT,
) {

	utils.PrintInfo("Announcing our node's presence")

	peerDiscoveryName := utils.GetConfig(
							"node.peerDiscoveryName", 
							"magicspace.ai",
						).(string)

	//for _,p := range rh.Peerstore().Peers(){
	//	rh.Peerstore().RemovePeer(p)
	//}

	routingDiscovery := drouting.NewRoutingDiscovery(dht)
	dutil.Advertise(ctx, routingDiscovery, peerDiscoveryName)
	
	var peerConnectFailCount = make(map[peer.ID]int64)
	maxPeerConnectRetries := utils.GetConfig("node.maxPeerConnectRetries", 10).(int64)
	peerDiscoveryScanInterval := utils.GetConfig("node.maxPeerConnectRetries", 30).(int64)


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
			
			//utils.PrintInfo("Peer has no addresses: %s", peerinfo.ID)

			if len(peerinfo.Addrs) == 0 {
				//utils.PrintInfo("Peer has no addresses: %s", peerinfo.ID.Pretty())
				rh.Peerstore().RemovePeer(peerinfo.ID)
				continue
			}

			if peerConnectFailCount[peerinfo.ID] > maxPeerConnectRetries {
				rh.Peerstore().RemovePeer(peerinfo.ID)
				continue		
			}

			// if already connected, skip it
			if rh.Network().Connectedness(peerinfo.ID) == network.Connected {
				continue
			}
			
			// if peer has no address lets skip it
			if len(peerinfo.Addrs) == 0 {
				continue
			}

			ethAddr, err := GetEthAddrFromPeer(peerinfo)

			if err != nil {
				utils.PrintError("Get eth address error, err=%v", err)
				continue
			}
		
			//err := rh.Connect(ctx, peerinfo)
			_, err = rh.Network().DialPeer(ctx, peerinfo.ID)

			if err != nil {
				fmt.Println("peer connection failed: ", peerinfo.ID.Pretty(), ", error:", err)
				peerConnectFailCount[peerinfo.ID] += 1
			} else {
				
				utils.PrintSuccess("Connected to: %s\n", peerinfo.ID.Pretty())
				fmt.Printf("Validator ETH Address: %s\n", ethAddr)

				peerConnectFailCount[peerinfo.ID] = 0

				ConnectedPeers[peerinfo.ID] = peerinfo

				rh.Peerstore().AddAddrs(
					peerinfo.ID, 
					peerinfo.Addrs,
					peerstore.TempAddrTTL,
				)

				
			}

			time.Sleep(100 * time.Nanosecond)
		}

		
		println()
		fmt.Printf("Total connected Peers %d \n", len(rh.Peerstore().Peers()))
		
		// if no peer was connected, lets try to re advertise our presence
		if len(ConnectedPeers) == 0 {
			utils.PrintInfo("re-announcing our node's presence...")
			dutil.Advertise(ctx, routingDiscovery, peerDiscoveryName)
		}

		time.Sleep(time.Duration(peerDiscoveryScanInterval))
	}

}	

// List to connected peers events 
func onPeerConnectionEvent(rh *rhost.RoutedHost)  {

	//sub, err := rh.EventBus().Subscribe(&event.EvtPeerConnectednessChanged{})

	//utils.HandleError(err, "failed to subscribe to connection events",true)

	///defer sub.Close()

	/*for {
		e, ok := <-sub.Out()

		if !ok {
			return
		}

		evtPeerChanged := e.(event.EvtPeerConnectednessChanged)

		evtPeerChanged

		/*if evtPeerChanged.Connectedness == network.Connected {
			fmt.Printf("New Peer Connected %s\n\n", evtPeerChanged.Peer.Pretty())
		} else {
			fmt.Printf("Peer Disconnected %s\n\n", evtPeerChanged.Peer.Pretty())
		}
	}*/
}

// disconnect all recent peers
func disconnectAllPeers(rh *rhost.RoutedHost) {
	for _,p := range rh.Network().Conns(){
		p.Close()
		rh.Network().ClosePeer(p.RemotePeer())
	}
}


/*/ PingNode ping the provided node
func PingNode(rh *rhost.RoutedHost, peerId peer.ID) {
	
	pctx, cancel := context.WithCancel(context.Background())
	
	defer cancel()

	ps := ping.PingService(rh)

	ts, err := ps.Ping(pctx, p)
	if err != nil {
		t.Fatal(err)
	}
}
*/

/*/ get a list of connected nodes 
func GetConnectedNodes(rhost rhost.RoutedHost) ([] peer.ID) {
	for _, pid := range rhost.Peerstore().PeersWithAddrs(){

	}
}*/

