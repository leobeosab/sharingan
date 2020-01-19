package helpers

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/leobeosab/sharingan/internal/models"
)

func PrintNmapScan(hosts ...models.Host) {
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)
	for _, h := range hosts {
		fmt.Printf("%s\n", h.Subdomain)
		fmt.Fprintln(w, "Port\tProtocol\tService\tState\t")
		if len(h.Ports) == 0 {
			fmt.Printf("No open ports found\n")
			w.Flush()
		}

		for _, p := range h.Ports {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t\n", p.ID, p.Protocol, p.ServiceName, p.State)
		}
		w.Flush()
	}
}
