package p2p

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
)


const MS_GLOBAL_TOPIC = "magicspace-global-topic"

// initialize the PubSub Engine
// this is resposible for sending and getting messages
// on the protocol
func initPubSub(
	ctx *context.Context, 
	rhost *rhost.RoutedHost,
) (
	*pubsub.PubSub, 
	*pubsub.Topic, 
	error,
) {

	ps, err := pubsub.NewGossipSub(*ctx, rhost)

	if err != nil {
		return nil, nil, fmt.Errorf("pubsub engine failed to initailize, err=%w", err)
	}

	globalTopic, err := ps.Join(MS_GLOBAL_TOPIC)

	if err != nil {
		return nil, nil, fmt.Errorf("pubsub failed to join topic, err=%w", err)
	}

	return ps, globalTopic, nil
}