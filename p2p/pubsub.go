package p2p

import (
	"context"
	"fmt"
	"path"

	lpubsub "github.com/libp2p/go-libp2p-pubsub"
	peer "github.com/libp2p/go-libp2p/core/peer"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/magicspace/supernode/utils"
)


const MS_GLOBAL_TOPIC = "magicspace-global-topic"

// initialize the PubSub Engine
// this is resposible for sending and getting messages
// on the protocol
func InitPubSub(
	ctx context.Context, 
	rhost *rhost.RoutedHost,
) (
	*lpubsub.PubSub, 
	*lpubsub.Topic, 
	error,
) {

	tracer, err :=  getTracer()

	if err != nil {
		return nil, nil, fmt.Errorf("pubsub tracer error, err=%w", err)
	}

	ps, err := lpubsub.NewGossipSub(ctx, rhost, lpubsub.WithEventTracer(tracer))

	if err != nil {
		return nil, nil, fmt.Errorf("pubsub engine failed to initailize, err=%w", err)
	}

	globalTopic, err := ps.Join(MS_GLOBAL_TOPIC)

	if err != nil {
		return nil, nil, fmt.Errorf("pubsub failed to join topic, err=%w", err)
	}

	sub, err := globalTopic.Subscribe()

	if err != nil {
		return nil, nil, fmt.Errorf("failed to subscribe to globalTopic, err=%w", err)
	}

	go handleSubscribe(rhost.ID(), ctx, sub)
	
	return ps, globalTopic, nil
}

func getTracer() (*lpubsub.PBTracer, error) {
	dataDir, err := utils.GetDataDir("store")

	if err != nil {
		return nil, err
	}

	filePath := path.Join(dataDir, "pubsub-trace.json")

	return  lpubsub.NewPBTracer(filePath)
}

func handleSubscribe(
	hostId peer.ID,
	ctx context.Context, 
	subscriber *lpubsub.Subscription,
) {
	for {
		
		msg, err := subscriber.Next(ctx)
		
		if err != nil {
			utils.HandleError(err, "subscription error", false)
			continue
		}

		if msg.ReceivedFrom == hostId {
			continue
		}

		fmt.Printf("got message: %s, from: %s\n", string(msg.Data), msg.ReceivedFrom.Pretty())

	}

}