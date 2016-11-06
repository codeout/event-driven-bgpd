package main

import (
	"./kafka"
	log "github.com/Sirupsen/logrus"
	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/config"
	gobgp "github.com/osrg/gobgp/server"
)

func main() {
	//log.SetLevel(log.DebugLevel)
	s := gobgp.NewBgpServer()
	go s.Serve()

	// start grpc api server. this is not mandatory
	// but you will be able to use `gobgp` cmd with this.
	g := api.NewGrpcServer(s, ":50051")
	go g.Serve()

	// global configuration
	global := &config.Global{
		Config: config.GlobalConfig{
			As:       65001, // Update accordingly here
			RouterId: "192.168.0.64",
		},
	}

	if err := s.Start(global); err != nil {
		log.Fatal(err)
	}

	// setup kafka
	client, err := kafka.NewProducer()
	if err != nil {
		log.Fatal(err)
	}

	// monitor new routes and peer state
	w := s.Watch(gobgp.WatchBestPath(), gobgp.WatchPeerState(false))
	for {
		select {
		case ev := <-w.Event():
			switch msg := ev.(type) {
			case *gobgp.WatchEventBestPath:
				client.SendPathListMessage(msg)
			case *gobgp.WatchEventPeerState:
				client.SendPeerStateMessage(msg)
			}
		}
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
