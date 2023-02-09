package p2p

import (
	"crypto/ecdsa"
	"fmt"
	"path"
	"strings"

	"github.com/decred/dcrd/dcrec/secp256k1"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	lcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/magicspace/supernode/utils"
)

type Identity struct {
	PrivateKey *ecdsa.PrivateKey
	Publickey  *ecdsa.PublicKey
}

func GetPrivateKey() (Identity, error){

	var identity Identity

	dataDir,_ := utils.GetAppDataDir()

	file := path.Join(dataDir, "identity.toml")

	appData, err := utils.GetAppData("identity")

	if err != nil {
		return identity, fmt.Errorf("failed to load identity file at %s, err=%w", file, err)
	}
	
	if !appData.IsSet("privateKey") {
		return identity, fmt.Errorf("privateKey field required at %s", file)
	}

	dataTrimed := strings.TrimSpace(appData.GetString("privateKey"))

	if len(dataTrimed) == 0 {
		return identity, fmt.Errorf("privateKey field empty at %s", file)
	}

	privKey, err := ethCrypto.HexToECDSA(dataTrimed)

	if err != nil {
		return identity,  fmt.Errorf("private ToECDSA error %w", err)
	}

	identity.PrivateKey = privKey
	identity.Publickey = &privKey.PublicKey

	return identity, nil 
}

// Convert libp2p pub key to eth address 
func PubKeyToEthAddress(data lcrypto.PubKey) (string, error) {
	
	dbytes, _ :=  data.Raw()
	k, err := secp256k1.ParsePubKey(dbytes)

	if err != nil {
		return "", err
	}

	address := ethCrypto.PubkeyToAddress(*k.ToECDSA()).Hex()


	return address, nil
}

// extract the public key an convert it eth address
func GetEthAddrFromPeer(p peer.AddrInfo) (string, error) {

	pubkey, err := p.ID.ExtractPublicKey()

	if err != nil {
		return "", fmt.Errorf("failed to extract publick from peerinfo, err=%w", err)
	}

	return PubKeyToEthAddress(pubkey)
}