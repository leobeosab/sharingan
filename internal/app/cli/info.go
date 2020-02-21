package cli

import (
	"fmt"

	"github.com/leobeosab/sharingan/internal/models"
	"github.com/leobeosab/sharingan/pkg/storage"
)

func PrintDomains(settings *models.ScanSettings) {
	e, p := storage.RetrieveOrCreateProgram(settings.Store, settings.Target)
	if !e {
		fmt.Errorf("Error no program %s found", settings.Target)
		return
	}

	for s, _ := range p.Hosts {
		fmt.Println(s)
	}
}
