package cli

import "github.com/timshannon/bolthold"

// Singleton settings object

var (
	s *settings
)

type settings struct {
	Target   string
	Threads  int
	NoPrompt bool
	Store    *bolthold.Store
}

type DNSSettings struct {
	SkipProbe   bool
	Rescan      bool
	ReplaceSubs bool
	RootDomain  string
	DNSWordlist string
}

// DirbSettings data for CLI dirb CLI command
type DirbSettings struct {
	Domain   string
	Wordlist string
}

func ScanSettings() *settings {
	if s != nil {
		return s
	}

	s = &settings{
		Threads: 25,
	}

	return s
}
