package beats

import (
	"log"
	"strings"

	"github.com/WatchBeam/etcbeat"
	etcd "github.com/coreos/etcd/client"
)

func Start(addresses string, key, value string) error {
	eclient, err := etcd.New(etcd.Config{
		Endpoints: strings.Split(addresses, ","),
		Transport: etcd.DefaultTransport,
	})
	if err != nil {
		return err
	}

	m := etcbeat.NewKeyMaintainer(etcd.NewKeysAPI(eclient), key, value)
	go m.Maintain()

	for err := range m.Errors() {
		log.Print(err)
	}

	return nil
}
