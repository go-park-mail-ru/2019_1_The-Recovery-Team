package database

import (
	"reflect"
	"testing"
)

type profile struct {
	ID       int
	Nickname string
	Record   int
	Win      int
	Loss     int
}

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

func TestCorrectFind(t *testing.T) {
	dbm, err := InitDatabaseManager("recoveryteam", "123456", "localhost", "sadislands")
	if err != nil {
		t.Errorf("test for OK Failed - can't connect to database")
	}
	expected := &profile{
		ID:       1,
		Nickname: "test",
		Record:   0,
		Win:      0,
		Loss:     0,
	}

	query := `SELECT id, nickname, record, win, loss FROM profile WHERE id = $1`
	result := &profile{}
	if err := dbm.Find(result, query); err != nil {
		t.Errorf("test for OK Failed - get error on correct data")
	}

	if *result != *expected {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, expected)
	}
}

func TestCorrectFindAll(t *testing.T) {
	dbm, err := InitDatabaseManager("recoveryteam", "123456", "localhost", "sadislands")
	if err != nil {
		t.Errorf("test for OK Failed - can't connect to database")
	}

	expected := []profile{
		{
			ID:       1,
			Nickname: "test",
			Record:   0,
			Win:      0,
			Loss:     0,
		},
		{
			ID:       2,
			Nickname: "Ivan",
			Record:   0,
			Win:      0,
			Loss:     0,
		},
	}

	query := `SELECT id, nickname, record, win, loss FROM profile`
	result := []profile{}
	if err := dbm.FindAll(&result, query); err != nil {
		t.Errorf("test for OK Failed - get error on correct data")
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, expected)
	}
}

func TestIncorrectFind(t *testing.T) {
	dbm := &Manager{}
	if err := dbm.Find(struct{}{}, ""); err == nil {
		t.Errorf("test for ERROR Failed - doesn't return error on incorrect data")
	}
}

func TestIncorrectFindAll(t *testing.T) {
	dbm := &Manager{}
	if err := dbm.FindAll([]struct{}{}, ""); err == nil {
		t.Errorf("test for ERROR Failed - doesn't return error on incorrect data")
	}
}
