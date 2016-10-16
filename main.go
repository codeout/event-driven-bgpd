package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/config"
	"github.com/osrg/gobgp/packet/bgp"
	gobgp "github.com/osrg/gobgp/server"
	table "github.com/osrg/gobgp/table"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)
	s := gobgp.NewBgpServer()
	go s.Serve()

	// start grpc api server. this is not mandatory
	// but you will be able to use `gobgp` cmd with this.
	g := api.NewGrpcServer(s, ":50051")
	go g.Serve()

	// global configuration
	global := &config.Global{
		Config: config.GlobalConfig{
			As:       64601,
			RouterId: "192.168.0.64",
		},
	}

	if err := s.Start(global); err != nil {
		log.Fatal(err)
	}

	// neighbor configuration
	n := &config.Neighbor{
		Config: config.NeighborConfig{
			NeighborAddress: "192.168.0.72",
			PeerAs:          64600,
		},
	}

	if err := s.AddNeighbor(n); err != nil {
		log.Fatal(err)
	}

	// add routes
	attrs := []bgp.PathAttributeInterface{
		bgp.NewPathAttributeOrigin(0),
		bgp.NewPathAttributeNextHop("192.168.0.64"),
		bgp.NewPathAttributeAsPath([]bgp.AsPathParamInterface{bgp.NewAs4PathParam(bgp.BGP_ASPATH_ATTR_TYPE_SEQ, []uint32{4000, 400000, 300000, 40001})}),
	}
	if _, err := s.AddPath("", []*table.Path{table.NewPath(nil, bgp.NewIPAddrPrefix(24, "10.0.0.0"), false, attrs, time.Now(), false)}); err != nil {
		log.Fatal(err)
	}

	// monitor new routes
	w := s.Watch(gobgp.WatchBestPath())
	for {
		select {
		case ev := <-w.Event():
			switch msg := ev.(type) {
			case *gobgp.WatchEventBestPath:
				for _, path := range msg.PathList {
					// do something useful
					bytes, _ := path.MarshalJSON()
					fmt.Println(string(bytes))
				}
			}
		}
	}
}
