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

func SaveProgram(s *bolthold.Store, p *models.Program) {
	err := s.Bolt().Update(func(tx *bbolt.Tx) error {
		err := s.TxInsert(tx, p.ProgramName, *p)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func UpdateProgram(s *bolthold.Store, p *models.Program) {
	err := s.Bolt().Update(func(tx *bbolt.Tx) error {
		err := s.TxUpdate(tx, p.ProgramName, *p)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func UpdateOrCreateProgram(s *bolthold.Store, p *models.Program) {
	e, _ := ProgramEntryExists(s, p.ProgramName)
	if e {
		UpdateProgram(s, p)
	} else {
		SaveProgram(s, p)
	}
}

// Bolthold Store reference and rootdomain string
func RetrieveProgram(s *bolthold.Store, p string) []models.Program {
	var results []models.Program

	err := s.Find(&results, bolthold.Where("ProgramName").Eq(p))

	if err != nil {
		log.Println("No previous scan found")
		results = []models.Program{models.Program{}}
	}

	return results
}

func ProgramEntryExists(s *bolthold.Store, p string) (bool, []models.Program) {
	results := RetrieveProgram(s, p)
	return len(results) > 0, results
}
