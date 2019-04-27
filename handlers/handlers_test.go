package handlers

import (
	"api/database"
	"api/models"
	"api/session"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProfileOK(t *testing.T) {
	var mockDb, _, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()
	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	dbm := &database.Manager{}
	sm := &session.Manager{}
	dbm.SetDB(sqlxDB)
	env := &models.Env{
		Dbm: dbm,
		Sm:  sm,
	}

	req, err := http.NewRequest("GET", "/profiles", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := GetProfile(env)
	handler.ServeHTTP(rr, req)
	expected := http.StatusInternalServerError
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expected)
	}
}

func TestPutProfile(t *testing.T) {
	var mockDb, _, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()
	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	dbm := &database.Manager{}
	sm := &session.Manager{}
	dbm.SetDB(sqlxDB)
	env := &models.Env{
		Dbm: dbm,
		Sm:  sm,
	}

	req, err := http.NewRequest("PUT", "/profiles", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := PutProfile(env)
	handler.ServeHTTP(rr, req)
	expected := http.StatusForbidden
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expected)
	}
}

func TestPostProfile(t *testing.T) {
	var mockDb, _, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()
	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	dbm := &database.Manager{}
	sm := &session.Manager{}
	dbm.SetDB(sqlxDB)
	env := &models.Env{
		Dbm: dbm,
		Sm:  sm,
	}

	req, err := http.NewRequest("POST", "/profiles", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := PostProfile(env)
	handler.ServeHTTP(rr, req)
	expected := http.StatusBadRequest
	if status := rr.Code; status != expected {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expected)
	}
}

func TestGetSession(t *testing.T) {

}

func TestPostSession(t *testing.T) {

}
