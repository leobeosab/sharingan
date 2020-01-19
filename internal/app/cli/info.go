package cli

import (
	"fmt"

	"github.com/leobeosab/sharingan/internal/models"
	"github.com/leobeosab/sharingan/pkg/storage"
)

func PrintDomains(settings *models.ScanSettings) {
	r := storage.RetrieveProgram(settings.Store, settings.Target)
	if len(r) == 0 {
		fmt.Errorf("Error no program %s found", settings.Target)
		return
	}
	p := r[0]

	for s, _ := range p.Hosts {
		fmt.Println(s)
	}
}
