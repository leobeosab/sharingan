package storage

import (
	"log"
	"os"

	"github.com/leobeosab/sharingan/internal/models"
	"github.com/timshannon/bolthold"
	"go.etcd.io/bbolt"
)

// Checks if there is a single path in the p slice and if there is it uses that for the db file
// Otherwise it uses a defaulting name in the users home directory
func OpenStore(p ...string) *bolthold.Store {
	var dbl string
	if len(p) > 0 && p[0] != "" {
		dbl = p[0]
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		dbl = home + "/.sharingan.bhdb"
	}

	store, err := bolthold.Open(dbl, 0600, nil)
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
		log.Println("No previous scan found")
		results = []models.ScanResults{models.ScanResults{}}
	}

	return results
}

func ScanEntryExists(s *bolthold.Store, d string) (bool, []models.ScanResults) {
	results := RetrieveScanResults(s, d)
	return len(results) > 0, results
}
