package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/WatchBeam/pingpong/go/beats"
	"github.com/WatchBeam/pingpong/go/server"
)

var (
	address    = flag.String("address", "127.0.0.1", "Host the pingpong server listens on.")
	port       = flag.Int("port", 6239, "Port the pingpong server will listen on. Clients will connect to publicHost:port to measure latency.")
	publicPing = flag.String("publicHost", "ws://127.0.0.1:6239", "This ingest's websocket address for ping/pong.")

	cert = flag.String("cert", "", "Path to the TLS certificate. If passed, TLS will be enabled.")
	key  = flag.String("key", "", "Path to the TLS key.")

	location     = flag.String("location", "", "Human-readable name where the server lives.")
	publicIngest = flag.String("publicIngest", "", "This ingest's public address to be stored in etcd.")

	etcdDir   = flag.String("etcdDir", "/ingests", "Directory in etcd for storing pingping hosts. If empty, etcd heartbeats will not be enabled.")
	etcdAddrs = flag.String("etcdAddr", "http://127.0.0.1:2379", "List of comma-delimited etcd hosts.")
	etcdTtl   = flag.Duration("etcdTtl", 10*time.Second, "Heartbeat TTL for this host in etcd.")
)

func main() {
	flag.Parse()

	if *etcdDir != "" {
		data, err := json.Marshal(struct {
			Ping     string `json:"ping"`
			Ingest   string `json:"ingest"`
			Location string `json:"location"`
		}{
			*publicPing,
			*publicIngest,
			*location,
		})

		if err != nil {
			panic(err)
		}

		go func() {
			log.Fatal(beats.Start(*etcdAddrs, *etcdDir+"/"+getID(), string(data)))
		}()
	}

	server.Listen(*address, *port, *cert, *key)
}

// Returns a random-ish ID for this server.
func getID() string {
	return fmt.Sprintf("%d-%d", os.Getpid(), rand.Int())
}
