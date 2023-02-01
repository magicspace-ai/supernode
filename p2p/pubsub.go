package p2p

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
)

type MagicMsg struct {
	peerId 		string,
	account 	string,
	topic 		string,
	data  		string
}

// initialize the PubSub Engine
// this is resposible for sending and getting messages
// on the protocol
func initPubSub(ctx *context.Context, host *host.Host) {

	ps, err := pubsub.NewGossipSub(ctx, host)

	if err != nil {
		return nil, ftm.Errorf("pubSub engine failed to initailize", err)
	}
	
}