package database

import (
	"api/models"
	"fmt"
	"testing"
)

func TestCorrectDBInit(t *testing.T) {
	dbm, err := InitDatabaseManager("../migrations/0_initial.sql")
	if dbm == nil || err != nil {
		t.Errorf("test for OK Failed - database manager isn't initialized with correct data")
	}

	data := models.ProfileCreate{
		Email:    "ivan@mail.ru",
		Nickname: "ivan",
		Password: "123456",
	}
	dbm.CreateProfile(&data)
	_, err = dbm.CreateProfile(&data)
	if err.Error() == "EmailAlreadyExists" {
		fmt.Println(err.Error())
	}
}
