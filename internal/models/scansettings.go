package models

import "github.com/timshannon/bolthold"

type ScanSettings struct {
	DNSWordlist string
	RootDomain  string
	Target      string
	Threads     int
	SkipProbe   bool
	Rescan      bool
	NoPrompt    bool
	ReplaceSubs bool
	Store       *bolthold.Store
}

type NMapSettings struct {
	Interactive string
}
