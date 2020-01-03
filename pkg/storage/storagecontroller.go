package storage

import (
	"log"

	"github.com/leobeosab/sharingan/internal/models"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

func OpenStore() *bolthold.Store {
	store, err := bolthold.Open("/tmp/boldt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return store
}

func SaveScan(s *bolthold.Store, scan *models.ScanResults) {
	err := s.Bolt().Update(func(tx *bbolt.Tx) error {
		err := s.TxInsert(tx, scan.RootDomain, *scan)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func UpdateScan(s *bolthold.Store, scan *models.ScanResults) {
	err := s.Bolt().Update(func(tx *bbolt.Tx) error {
		err := s.TxUpdate(tx, scan.RootDomain, *scan)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

// Bolthold Store reference and rootdomain string
func RetrieveScanResults(s *bolthold.Store, d string) []models.ScanResults {
	var results []models.ScanResults

	err := s.Find(&results, bolthold.Where("RootDomain").Eq(d))

	if err != nil {
		log.Fatal(err)
	}

	return results
}
