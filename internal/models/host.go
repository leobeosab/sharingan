package models

// Shouldn't rely on IP can change based on load balancers and plenty of other variables
type Host struct {
	Subdomain   string
	Ports       []Port
	Http, Https bool
}

type Port struct {
	ID          uint16
	Protocol    string
	State       string
	ServiceName string
}
