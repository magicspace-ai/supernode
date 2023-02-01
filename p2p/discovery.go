package p2p

import (
	"context"
	"fmt"
	"sync"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/magicspace/supernode/utils"
	ma "github.com/multiformats/go-multiaddr"
)

// initialize the DHT engine, if it fails,
//a nil value with an eror will be returned
func initDHT(
	ctx context.Context, 
	hostPtr *host.Host,
) (
	*dht.IpfsDHT, 
	*rhost.RoutedHost, 
	error,
) {
	
	bootNodesList := utils.GetConfigs().GetStringSlice("node.bootNodes")

	var bootNodes []ma.Multiaddr

	if len(bootNodesList) == 0 {
		bootNodes = dht.DefaultBootstrapPeers
	} else {
		for _, addr := range bootNodesList {
			maAddr, _ := ma.NewMultiaddr(addr)
			bootNodes = append(bootNodes, maAddr)
		}
	}
	
	host := *hostPtr

	// Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
	dstore := dsync.MutexWrap(ds.NewMapDatastore())


	kDHT := dht.NewDHT(ctx, host, dstore)

	// Make the routed host
	routedHost := rhost.Wrap(host, kDHT)

	//if err != nil {
	//	return nil,fmt.Errorf("Failed to start DHT engine, err=%w", err)
	//}

	if err := kDHT.Bootstrap(ctx); err != nil {
		return nil, nil, fmt.Errorf("bootstraping DHT failed, err=%w", err)
	}

	var wg sync.WaitGroup
	
	for _, nodeAddr := range bootNodes {
		
		bnode, _ := peer.AddrInfoFromP2pAddr(nodeAddr)

		wg.Add(1)

		go func() {
			
			defer wg.Done()
			err := host.Connect(ctx, *bnode);
			
			if err != nil {
				fmt.Sprintf("Failed to connect to bootnode %s", bnode.ID)
			} else {
				fmt.Sprintf("Connected to bootnode %s", bnode.ID)
			}	
		}()
	}

	wg.Wait()

	return kDHT, routedHost, nil
}

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

