package models

import "github.com/timshannon/bolthold"

type ScanSettings struct {
	DNSWordlist string
	Target      string
	SkipProbe   bool
	Store       *bolthold.Store
}
