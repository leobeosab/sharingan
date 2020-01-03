package storage

import (
	"log"

	"github.com/timshannon/bolthold"
)

func OpenStore() *bolthold.Store {
	store, err := bolthold.Open("/tmp/boldt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return store
}
