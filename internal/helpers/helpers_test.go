package helpers

import (
	"bufio"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// e = expected, a = actual

func TestGetKeysFromMap(t *testing.T) {
	t.Logf("Testing helpers.GetKeysFromMap")

	m := map[string]string{
		"192.168.1.1": "asus.routerlogin.com",
		"127.0.0.1":   "localhost",
	}

	e := []string{"192.168.1.1", "127.0.0.1"}
	a := GetKeysFromMap(&m)

	if !cmp.Equal(e, a) {
		t.Errorf("Error slices don't match, expected: %v, got %v", e, a)
	}
}

func TestGetNumberOfLinesInFile(t *testing.T) {
	t.Logf("Testing helpers.GetNumberOfLinesInFile")

	e := 4

	tmp, err := ioutil.TempFile(os.TempDir(), "sharingantesting-helpers-")
	if err != nil {
		t.Fatalf("Error creating temporary file")
	}
	defer tmp.Close()
	defer os.Remove(tmp.Name())

	w := bufio.NewWriter(tmp)
	_, err = w.WriteString("one\ntwo\nthree\nfour")
	if err != nil {
		t.Fatalf("Error writing to temporary file")
	}
	w.Flush()

	a := GetNumberOfLinesInFile(tmp)
	if a != e {
		t.Fatalf("Error expected: %s actual: %s", strconv.Itoa(e), strconv.Itoa(a))
	}
}
