package database

import "testing"

func TestCorrectDBInit(t *testing.T) {
	dbm, err := InitDatabaseManager("recoveryteam", "123456", "localhost", "sadislands")
	if dbm == nil || err != nil {
		t.Errorf("test for OK Failed - database manager isn't initialized with correct data")
	}
}

func TestIncorrectDBInit(t *testing.T) {
	dbm, err := InitDatabaseManager("", "", "", "")
	if dbm != nil || err == nil {
		t.Errorf("test for ERROR Failed - doesn't return error on database manager initialization with incorrect data")
	}
}

func TestCorrectDBManagerGet(t *testing.T) {
	dbm, _ := InitDatabaseManager("recoveryteam", "123456", "localhost", "sadislands")
	db, err := dbm.DB()
	if err != nil || db == nil {
		t.Errorf("test for OK Failed - doesn't return initialized database instance")
	}
}

func TestIncorrectDBManagerGet(t *testing.T) {
	dbm := &Manager{}
	db, err := dbm.DB()
	if db != nil || err == nil {
		t.Errorf("test for ERROR Failed - doesn't return error on getting uninitialized database instance")
	}
}

func TestCorrectDBClose(t *testing.T) {
	dbm, _ := InitDatabaseManager("recoveryteam", "123456", "localhost", "sadislands")
	err := dbm.Close()
	if err != nil {
		t.Errorf("test for OK Failed - doesn't correctly close database connection")
	}
}
func TestIncorrectDBClose(t *testing.T) {
	dbm := &Manager{}
	err := dbm.Close()
	if err == nil {
		t.Errorf("test for ERROR Failed - doesn't return error incorrect database connection close")
	}
}
