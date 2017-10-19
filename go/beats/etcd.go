package beats

import (
	"log"
	"strings"

	etcd "github.com/coreos/etcd/client"
	"github.com/mixer/etcbeat"
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
