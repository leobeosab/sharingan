package models

import "time"

type Scan struct {
	RootDomain      string
	Hosts           []Host
	DateLastScanned time.Time
}
