package main

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-discovery"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

func setupMDNS(ctx context.Context, h host.Host, serviceTag string) error {
	mdnsService, err := mdns.NewMdnsService(ctx, h, time.Minute, serviceTag)
	if err != nil {
		return err
	}

	disc := discovery.NewRoutingDiscovery(mdnsService)
	discovery.Advertise(ctx, disc, serviceTag)

	mdnsService.RegisterNotifee(&mdnsNotifee{h: h})
	return nil
}

// mdnsNotifee implements the Notifee interface for mDNS discovery
type mdnsNotifee struct {
	h host.Host
}

func (n *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("Discovered new peer: %s\n", pi.ID.Pretty())
	n.h.Peerstore().AddAddrs(pi.ID, pi.Addrs, peer.AddressTTL)
}