package storage

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/leobeosab/sharingan/internal/models"
	"github.com/timshannon/bolthold"
)

func TestOpenStore(t *testing.T) {
	t.Logf("Testing OpenStore no custom path function \n\n")
	store := OpenStore()
	defer store.Close()
}

func TestStoreAndRetrieve(t *testing.T) {
	t.Logf("Testing Save & Retrieve Scan with custom db path")

	expected := NewBasicScanResults()
	s := CreateBasicStoreAndEntry(expected, t)

	ret := RetrieveScanResults(s, expected.RootDomain)
	actual := ret[0]

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Compare scan results mismatch (-want +got):\n%s", diff)
	}
}

func TestUpdate(t *testing.T) {
	t.Logf("Testing Update scan results")

	expected := NewBasicScanResults()
	s := CreateBasicStoreAndEntry(expected, t)

	expected.Hosts = append(expected.Hosts,
		models.Host{
			IP:         "192.168.1.1",
			Subdomains: []string{"router.asus"},
			OpenPorts:  []int{22},
			Http:       false,
		},
	)

	UpdateScan(s, &expected)

	ret := RetrieveScanResults(s, expected.RootDomain)
	actual := ret[0]

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Compare scan results mismatch (-want +got):\n%s", diff)
	}
}

func CreateBasicStoreAndEntry(m models.ScanResults, t *testing.T) *bolthold.Store {
	// Create tmp file so we don't get same id errors
	tmp, err := ioutil.TempFile(os.TempDir(), "sharingantesting-")
	if err != nil {
		t.Fatalf("Error creating temporary file")
	}
	defer os.Remove(tmp.Name())

	s := OpenStore(tmp.Name())
	SaveScan(s, &m)

	return s
}

func NewBasicScanResults() models.ScanResults {
	return models.ScanResults{
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
}
