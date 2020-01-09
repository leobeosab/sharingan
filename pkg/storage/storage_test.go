package storage

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/leobeosab/sharingan/internal/models"
)

func TestOpenStore(t *testing.T) {
	t.Logf("Testing OpenStore no custom path function \n\n")
	store := OpenStore()
	defer store.Close()
}

func TestStoreAndRetrieve(t *testing.T) {
	t.Logf("Testing Save & Retrieve Scan with custom db path")

	expected := models.ScanResults{
		RootDomain: "bestdomain.com",
		Hosts: []models.Host{
			models.Host{
				IP:         "127.0.0.1",
				Subdomains: []string{"self.home.com", "ride.fast", "eat.ass"},
				OpenPorts:  []int{80, 443},
				Http:       true,
			},
		},
		DateLastScanned: time.Now(),
	}

	// Create tmp file so we don't get same id errors
	tmp, err := ioutil.TempFile(os.TempDir(), "sharingantesting-")
	if err != nil {
		t.Fatalf("Error creating temporary file")
	}
	defer os.Remove(tmp.Name())

	s := OpenStore(tmp.Name())
	SaveScan(s, &expected)

	ret := RetrieveScanResults(s, expected.RootDomain)
	actual := ret[0]

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Compare scan results mismatch (-want +got):\n%s", diff)
	}
}
