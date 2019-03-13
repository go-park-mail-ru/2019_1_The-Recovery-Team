package database

import (
	"fmt"
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
	dbm, err := InitDatabaseManager("recoveryteam", "123456", "localhost:5432", "test", "../migrations/test", true)
	if dbm == nil || err != nil {
		t.Errorf("test for OK Failed - database manager isn't initialized with correct data")
	}
}

func TestIncorrectDBInit(t *testing.T) {
	dbm, err := InitDatabaseManager("", "", "", "", "", true)
	if dbm != nil || err == nil {
		t.Errorf("test for ERROR Failed - doesn't return error on database manager initialization with incorrect data")
	}
}

func TestCorrectDBManagerGet(t *testing.T) {
	dbm, _ := InitDatabaseManager("recoveryteam", "123456", "localhost:5432", "test", "../migrations/test", true)
	db := dbm.DB()
	if db == nil {
		t.Errorf("test for OK Failed - doesn't return initialized database instance")
	}
}

func TestCorrectDBClose(t *testing.T) {
	dbm, _ := InitDatabaseManager("recoveryteam", "123456", "localhost:5432", "test", "../migrations/test", true)
	err := dbm.Close()
	if err != nil {
		t.Errorf("test for OK Failed - doesn't correctly close database connection")
	}
}
func TestIncorrectDBClose(t *testing.T) {
	dbm := &Manager{}
	err := dbm.Close()
	if err != nil {
		t.Errorf("test for ERROR Failed - doesn't return error incorrect database connection close")
	}
}

func TestCorrectFind(t *testing.T) {
	dbm, err := InitDatabaseManager("recoveryteam", "123456", "localhost:5432", "test", "../migrations/test", true)
	if err != nil {
		t.Errorf("test for OK Failed - can't connect to database")
		return
	}
	expected := &profile{
		ID:       100,
		Nickname: "test",
	}

	query := `SELECT id, nickname, record, win, loss FROM profile WHERE id = $1`
	result := &profile{}
	if err := dbm.Find(result, query, expected.ID); err != nil {
		fmt.Println(err)
		t.Errorf("test for OK Failed - get error on correct data")
		return
	}

	if *result != *expected {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, expected)
	}
}

func TestCorrectFindAll(t *testing.T) {
	dbm, err := InitDatabaseManager("recoveryteam", "123456", "localhost:5432", "test", "../migrations/test", true)
	if err != nil {
		t.Errorf("test for OK Failed - can't connect to database")
		return
	}

	expected := []profile{
		{
			ID:       100,
			Nickname: "test",
		},
		{
			ID:       101,
			Nickname: "test3",
		},
		{
			ID:       102,
			Nickname: "test5",
		},
	}

	query := `SELECT id, nickname, record, win, loss FROM profile`
	result := []profile{}
	if err := dbm.FindAll(&result, query); err != nil {
		t.Errorf("test for OK Failed - get error on correct data")
		return
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, expected)
	}
}
