package p2p

import (
	"crypto/rand"
	"fmt"

	"github.com/magicspace/supernode/utils"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	host "github.com/libp2p/go-libp2p/core/host"
	protocol "github.com/libp2p/go-libp2p/core/protocol"
)

/**
 * create node
 */
func MakeNode() (*host.Host, error) {

	hostIp := utils.GetConfig("node.host", "0.0.0.0").(string)
	port   := utils.GetConfig("node.port", 60_000).(int64)
	protocolId := utils.GetConfig("node.protocolId", "magicspace://").(string)

	identity, err := utils.GetAppData("identity")

	if err != nil {
		return nil, err
	}

	
	var priv crypto.PrivKey

	if !(identity.IsSet("privateKey") && identity.Get("privateKey").(string) == "") {

		priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)

		if err != nil {
			return nil, err
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
			return nil, err
		}

		priv, _, err = crypto.KeyPairFromStdKey(privBytes)

		if err != nil{
			return nil, err
		}	
			
	}
	
	
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", hostIp, port)),
		libp2p.Identity(priv),
	}
	
	
	hostInfo, err := libp2p.New(opts...)

	
	if err != nil {
		return nil, fmt.Errorf("failed to initialized libp2p node err=%s", err)
	}

	fmt.Printf("Node created with ID: %s\n", hostInfo.Addrs()[0])

	hostInfo.SetStreamHandler(protocol.ID(protocolId), streamHandler)

	return &hostInfo, nil
}
