package models

import "time"

type ScanResults struct {
	RootDomain      string
	Hosts           []Host
	DateLastScanned time.Time
}
