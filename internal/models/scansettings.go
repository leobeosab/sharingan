package models

import "github.com/timshannon/bolthold"

type ScanSettings struct {
	DNSWordlist string
	Target      string
	SkipProbe   bool
	Rescan      bool
	Store       *bolthold.Store
}

type NMapSettings struct {
	Interactive string
}
