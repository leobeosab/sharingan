package models

import "time"

// Bug bounty program

type Program struct {
	ProgramName     string    // Name of program
	Hosts           []Host    // Actual hosts and subs listed with hosts
	Subdomains      []string  // Full list of subdomains
	DateLastScanned time.Time // Last date
}
