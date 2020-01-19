package models

import "time"

// Bug bounty program

type Program struct {
	ProgramName     string          // Name of program
	Hosts           map[string]Host // Actual hosts and subs listed with hosts
	DateLastScanned time.Time       // Last date
}
