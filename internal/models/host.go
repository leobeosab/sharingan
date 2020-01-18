package models

type Host struct {
	IP          string
	Subdomains  []string
	Ports       []Port
	Http, Https bool
}

type Port struct {
	ID          uint16
	Protocol    string
	State       string
	ServiceName string
}
