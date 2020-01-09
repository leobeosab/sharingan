package helpers

import (
	"reflect"
	"testing"
)

// e = expected, a = actual

func TestGetKeysFromMap(t *testing.T) {
	m := map[string]string{
		"192.168.1.1": "asus.routerlogin.com",
		"127.0.0.1":   "localhost",
	}

	e := []string{"192.168.1.1", "127.0.0.1"}
	a := GetKeysFromMap(&m)

	if !reflect.DeepEqual(e, a) {
		t.Errorf("Error slices don't match, expected: %v, got %v", e, a)
	}
}
