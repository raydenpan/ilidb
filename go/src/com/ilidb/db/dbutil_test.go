package db

import (
	"testing"
)

func TestDBFetchUserDB(t *testing.T) {
	tCollection := getDatabaseCollection("ilidb_users")
	if nil == tCollection {
		t.Fail()
	}
}
