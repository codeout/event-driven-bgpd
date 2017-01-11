# event-driven-bgpd

Sandbox of Event Driven BGPd.

## BGP KONAMI Command

A event driven BGPd which accepts KONAMI command (↑↑↓↓←→←→BA) in BGP, and then advertises full route to neighbors while it advertises default route only in normal operation.

### Prerequisites

* Apache kafka
* GoBGP
* Ruby with gRPC environment

### How to play

In OSX,

#### Start kafka

1. Start Zookeeper and kafka
   ```zsh
   zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties & kafka-server-start /usr/local/etc/kafka/server.properties
   ```

2. Create a Topic
   ```zsh
   kafka-topics --create --zookeeper localhost:2181 --topic edb --partitions 1 --replication-factor 1
   ```

#### BGPd

```zsh
git clone https://github.com/codeout/event-driven-bgpd.git
cd event-driven-bgpd
```

Let's say the sandbox enviroment looks like this:

```
+-----------------------+      eBGP      +--------+
| KONAMI Command Router +----------------+ Client |
+-----------------------+                +--------+
   192.168.0.64 / AS65001                192.168.0.71 / AS65000
```

1. Build and Start bgpd
   ```zsh
   go get github.com/Shopify/sarama
   cd bgpd
   go build .   # requires GoBGP
   bgpd &       # assumes no password prompt
   ```

2. Download MRT and configure bgpd

   Download MRT table dump from [Route View Archive](http://archive.routeviews.org/) for example.
   ```zsh
   curl -O http://archive.routeviews.org/route-views.wide/bgpdata/2017.01/RIBS/rib.20170101.0200.bz2
   bunzip2 rib.20170101.0200.bz2
   ./config.sh rib.20170101.0200  # this will take a while
   ```

#### KONAMI Command subscriber

Install gRPC and Ruby libraries in advance. (See [GoBGP doc](https://github.com/osrg/gobgp/blob/master/docs/sources/grpc-client.md#ruby) for detail)

1. Generate Stub Code
   ```zsh
   cd ../konami
   GOBGP_API=$GOPATH/src/github.com/osrg/gobgp/api
   grpc_tools_ruby_protoc -I $GOBGP_API --ruby_out=. --grpc_out=. $GOBGP_API/gobgp.proto
   ```

   You can find ```gobgp_pb.rb``` and ```gobgp_services_pb.rb``` now.

2. Start Consumer
   ```zsh
   gem install ruby-kafla grpc
   ruby subscriber.rb &
   ```

#### BGP Client

In other system,

```
git clone https://github.com/codeout/event-driven-bgpd.git
cd event-driven-bgpd/konami_client
```

1. Start Client
   ```zsh
   sudo gobgpd -f gobgpd.conf
   ```

   You can see a route received after BGP peer is established.

   ```zsh
   gobgp neighbor
   
   Peer            AS  Up/Down State       |#Advertised Received Accepted
   192.168.0.64 65001 00:00:35 Establ      |          0        1        1
   ```

2. Send KONAMI Command

   You will receive full route immediately after send sequencial BGP updates and reset the peer.
   ```zsh
   ./send_command.sh
   gobgp neighbor

   Peer            AS  Up/Down State       |#Advertised Received Accepted
192.168.0.64 65001 00:00:34 Establ      |          0   612259   612259
   ```

### Message Example

```json
{"type":"best_path","value":{"nlri":{"prefix":"10.0.0.0/8"},"attrs":[{"type":1,"value":2},{"type":2,"as_paths":[{"segment_type":2,"num":1,"asns":[65000]}]},{"type":3,"nexthop":"192.168.0.71"}],"age":1476773724,"source-id":"192.168.0.71","neighbor-ip":"192.168.0.71"}}
```
