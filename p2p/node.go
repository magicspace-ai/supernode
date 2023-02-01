package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"

	"github.com/magicspace/supernode/utils"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	host "github.com/libp2p/go-libp2p/core/host"
	protocol "github.com/libp2p/go-libp2p/core/protocol"

	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	ma "github.com/multiformats/go-multiaddr"
)

/**
 * create node
 */
func MakeNode(ctx context.Context) (
	*host.Host, 
	*dht.IpfsDHT, 
	*rhost.RoutedHost, 
	error,
) {

	hostIp := utils.GetConfig("node.host", "0.0.0.0").(string)
	port   := utils.GetConfig("node.port", 60_000).(int64)
	protocolId := utils.GetConfig("node.protocolId", "magicspace://").(string)

	identity, err := utils.GetAppData("identity")

	if err != nil {
		return nil, nil, nil, err
	}

	
	var priv crypto.PrivKey

	if !(identity.IsSet("privateKey") && identity.Get("privateKey").(string) == "") {

		priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

		if err != nil {
			return nil, nil, nil, err
		}

		privkeyBytes, _ := priv.Raw()
		pubKeyBytes, _ := priv.GetPublic().Raw()

		dataToSave := map[string]interface{}{
			"privateKey": utils.ToHex(privkeyBytes),
			"publicKey":  utils.ToHex(pubKeyBytes),
		}

		utils.SaveAppData("identity", dataToSave)

	} else {

		privBytes, err := utils.FromHex(identity.Get("privateKey").(string))

		if err != nil{
			return nil, nil, nil, err
		}

		priv, _, err = crypto.KeyPairFromStdKey(privBytes)

		if err != nil{
			return nil, nil, nil, err
		}	
			
	}
	
	
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", hostIp, port)),
		libp2p.Identity(priv),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
	}
	
	
	hostInfo, err := libp2p.New(opts...)

	
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to initialized libp2p node err=%s", err)
	}

	fmt.Printf("Node created with ID: %s\n", hostInfo.Addrs()[0])

	hostInfo.SetStreamHandler(protocol.ID(protocolId), streamHandler)

	kDHT, rhost, err := initDHT(ctx, &hostInfo)

	if err != nil {
		return nil, nil, nil, err 
	}

	return &hostInfo, kDHT, rhost, nil
}

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
