package cli

import (
	"fmt"

	"github.com/leobeosab/sharingan/pkg/storage"
)

func PrintDomains() {
	e, p := storage.RetrieveOrCreateProgram(ScanSettings().Store, ScanSettings().Target)
	if !e {
		fmt.Errorf("Error no program %s found", ScanSettings().Target)
		return
	}

	for s, _ := range p.Hosts {
		fmt.Println(s)
	}
}
