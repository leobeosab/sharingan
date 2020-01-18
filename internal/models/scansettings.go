package models

import "github.com/timshannon/bolthold"

type ScanSettings struct {
	DNSWordlist string
	RootDomain  string
	Target      string
	SkipProbe   bool
	Rescan      bool
	ReplaceSubs bool
	Store       *bolthold.Store
}

type NMapSettings struct {
	Interactive string
}
