package storage

import (
	"io/ioutil"
	"os"
	"testing"

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

	expected := NewBasicProgram()
	s := CreateBasicStoreAndEntry(expected, t)

	ret := RetrieveProgram(s, expected.ProgramName)
	actual := ret[0]

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Compare scan results mismatch (-want +got):\n%s", diff)
	}
}

func TestUpdate(t *testing.T) {
	t.Logf("Testing Update scan results")

	expected := NewBasicProgram()
	s := CreateBasicStoreAndEntry(expected, t)

	expected.Hosts["testsubdomain"] =
		models.Host{
			Subdomain: "testsubdomain",
		}

	UpdateProgram(s, &expected)

	ret := RetrieveProgram(s, expected.ProgramName)
	actual := ret[0]

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Compare scan results mismatch (-want +got):\n%s", diff)
	}
}

func CreateBasicStoreAndEntry(m models.Program, t *testing.T) *bolthold.Store {
	// Create tmp file so we don't get same id errors
	tmp, err := ioutil.TempFile(os.TempDir(), "sharingantesting-")
	if err != nil {
		t.Fatalf("Error creating temporary file")
	}
	defer os.Remove(tmp.Name())

	s := OpenStore(tmp.Name())
	SaveProgram(s, &m)

	return s
}

func NewBasicProgram() models.Program {
	return models.Program{
		ProgramName: "test",
		Hosts:       map[string]models.Host{"testhost": models.Host{}},
	}
}
