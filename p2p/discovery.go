package p2p

import (
	"context"
	"fmt"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

//Start the DHT engine, if it fails, a nil value with an eror will be returned
func startDHT(ctx context.Context, hostPtr *host.Host) (*dht.IpfsDHT, error) {
	
	kDHT, err := dht.New(ctx, hostPtr)

	if err != nil {
		return nil,fmt.Errorf("Failed to start DHT engine, err=%w", err)
	}

	return kDHT, nil
}


// Discover peers using the boostrap method
// This helps other nodes on the network find themseleves
func discoverPeers(ctx context.Context, hostPtr *host.Host){
	
}

