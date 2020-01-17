package models

import "time"

// Bug bounty program

type Program struct {
	ProgramName     string
	Hosts           []Host
	DateLastScanned time.Time
}
