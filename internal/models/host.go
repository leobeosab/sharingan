package models

type Host struct {
	IP          string
	Subdomains  []string
	OpenPorts   []int
	Http, https bool
}
