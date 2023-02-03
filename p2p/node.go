package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"

	"github.com/magicspace/supernode/core"
	"github.com/magicspace/supernode/utils"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	host "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	protocol "github.com/libp2p/go-libp2p/core/protocol"

	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	dquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	ma "github.com/multiformats/go-multiaddr"
)

//var logger = log.Logger("p2p")

/**
 * create node
 */
func MakeNode( ctx context.Context ) (
	*rhost.RoutedHost, 
	*dht.IpfsDHT, 
	error,
) {

	resourceManager, err := core.GetResourceMgr()

	if err != nil {
		return nil, nil, err
	}

	hostIp := utils.GetConfig("node.host", "0.0.0.0").(string)
	port   := utils.GetConfig("node.port", 60_000).(int64)
	protocolId := utils.GetConfig("node.protocolId", "magicspace://").(string)

	identity, err := utils.GetAppData("identity")

	if err != nil {
		return nil, nil, err
	}

	var priv crypto.PrivKey

	if !(identity.IsSet("privateKey") && identity.Get("privateKey").(string) == "") {

		priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

		if err != nil {
			return nil, nil, err
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
			return nil, nil, err
		}

		priv, _, err = crypto.KeyPairFromStdKey(privBytes)

		if err != nil{
			return nil, nil, err
		}	
			
	}
		
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/%s/tcp/%d", hostIp, port),
			fmt.Sprintf("/ip4/%s/tcp/%d/ws", hostIp, port),
			fmt.Sprintf("/ip4/%s/tcp/%d/quic", hostIp, port),
		),
		libp2p.Identity(priv),
		libp2p.NoTransports,

		//libp2p.EnableAutoRelay(),
		libp2p.EnableRelayService(),
		libp2p.EnableNATService(),
		libp2p.ForceReachabilityPublic(),
		libp2p.EnableRelay(),
		libp2p.NATPortMap(),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(websocket.New),
		libp2p.Transport(dquic.NewTransport),
		libp2p.DefaultMuxers,
		libp2p.ResourceManager(resourceManager),
	}
	
	
	hostInfo, err := libp2p.New(opts...)

	
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialized libp2p node err=%s", err)
	}

	utils.PrintSuccess("Node created with Address: %s\n", hostInfo.Addrs()[0])

	hostInfo.SetStreamHandler(protocol.ID(protocolId), streamHandler)

	kDHT, rhost, err := initDHT(ctx, &hostInfo)

	if err != nil {
		return nil, nil, err 
	}

	utils.PrintSuccess("Node created with ID: %s\n", rhost.ID().Pretty())

	return  rhost, kDHT, nil
}

// initialize the DHT engine, if it fails,
//a nil value with an eror will be returned
func  initDHT(
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
	//dstore := dsync.MutexWrap(ds.NewMapDatastore())

	var dhtBootPeers []peer.AddrInfo

	// convert the boot nodes to dht peer addr
	for _,bn := range bootNodes {
		bnAddrInfo,_ := peer.AddrInfoFromP2pAddr(bn)
		dhtBootPeers = append(dhtBootPeers, *bnAddrInfo)
	}

	dhtOpts := []dht.Option{
		//dht.Mode(dht.ModeAutoServer),
		dht.BootstrapPeers(dhtBootPeers...),
	}

	kDHT, err := dht.New(ctx, host, dhtOpts...)

	if err != nil {
		return nil, nil, fmt.Errorf("kDHT engine failed to initialize, err=%w", err)
	}

	// Make the routed host
	routedHost := rhost.Wrap(host, kDHT)

	if err := kDHT.Bootstrap(ctx); err != nil {
		return nil, nil, fmt.Errorf("bootstraping DHT failed, err=%w", err)
	}

	var wg sync.WaitGroup
	
	for _, bn := range dhtBootPeers {
		
		wg.Add(1)

		go func() {
			
			defer wg.Done()
			err := host.Connect(ctx, bn);
			
			if err != nil {
				fmt.Printf("failed to connect to bootnode %s\n", bn.ID.Pretty())
			} else {
				fmt.Printf("connected to bootnode %s\n", bn.ID.Pretty())
			}	
		}()
	}

	fmt.Println("")
	wg.Wait()

	return kDHT, routedHost, nil
}
