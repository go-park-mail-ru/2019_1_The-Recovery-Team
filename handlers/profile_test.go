package handlers

import (
	"api/database"
	"api/middleware"
	"api/models"
	"api/session"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
)

type TestCase struct {
	Body       io.Reader
	Response   string
	StatusCode int
}

type TestCaseGetProfile struct {
	ProfileID uint64
	TestCase
}

type TestCaseCreateProfile struct {
	ProfileRegistration models.ProfileRegistration
	TestCase
}

type TestCasePutProfile struct {
	ProfileID uint64
	TestCase
}

type TestCaseGetProfiles struct {
	Email      string
	Nickname   string
	StatusCode int
}

var (
	dbm, _ = database.InitDatabaseManager("recoveryteam", "123456", "localhost:5432", "test", "../migrations/test", true)
	sm, _  = session.InitSessionManager("", "", "localhost:6379")
)

var env = models.Env{
	Dbm: dbm,
	Sm:  sm,
}

const baseUrl = "http://127.0.0.1:8080/api/v1"

func TestGetProfile(t *testing.T) {
	profile := &models.Profile{
		ID: 100,
		ProfileInfo: models.ProfileInfo{
			Email:    "test@mail.ru",
			Nickname: "test",
		},
	}
	profile.Email = ""
	expected, _ := easyjson.Marshal(profile)

	cases := []TestCaseGetProfile{
		{
			ProfileID: profile.ID,
			TestCase: TestCase{
				Response:   string(expected),
				StatusCode: http.StatusOK,
			},
		},
		{
			ProfileID: 0,
			TestCase: TestCase{
				Response:   "",
				StatusCode: http.StatusNotFound,
			},
		},
	}

	for number, item := range cases {
		u := baseUrl + "/profiles/"
		vars := map[string]string{
			"id": strconv.FormatUint(item.ProfileID, 10),
		}
		req := httptest.NewRequest("GET", u, nil)
		req = mux.SetURLVars(req, vars)
		w := httptest.NewRecorder()
		GetProfile(&env)(w, req)
		expected := item.StatusCode
		if w.Code != expected {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				number, w.Code, expected)
			return
		}

		resp := w.Result()

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		resp.Body.Close()
		body := string(b)
		if body != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				number, body, item.Response)
			return
		}
	}

}

func TestPostProfile(t *testing.T) {
	expected1 := models.Profile{
		ProfileInfo: models.ProfileInfo{
			Email:    "test1@mail.ru",
			Nickname: "test1",
		},
	}

	received1 := models.ProfileRegistration{
		ProfileLogin: models.ProfileLogin{
			Email:    "test1@mail.ru",
			Password: "123",
		},
		Nickname: "test1",
	}

	received2 := models.ProfileRegistration{}

	expected3 := models.HandlerError{
		Description: "email already exists",
	}

	received3 := models.ProfileRegistration{
		ProfileLogin: models.ProfileLogin{
			Email:    "test3@mail.ru",
			Password: "123",
		},
		Nickname: "test100",
	}

	expected4 := models.HandlerError{
		Description: "nickname already exists",
	}

	received4 := models.ProfileRegistration{
		ProfileLogin: models.ProfileLogin{
			Email:    "test100@mail.ru",
			Password: "123",
		},
		Nickname: "test3",
	}

	e1, _ := easyjson.Marshal(expected1)
	e3, _ := easyjson.Marshal(expected3)
	e4, _ := easyjson.Marshal(expected4)

	cases := []TestCaseCreateProfile{
		{
			ProfileRegistration: received1,
			TestCase: TestCase{
				Response:   string(e1),
				StatusCode: http.StatusOK,
			},
		},
		{
			ProfileRegistration: received2,
			TestCase: TestCase{
				Response:   "",
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			ProfileRegistration: received3,
			TestCase: TestCase{
				Response:   string(e3),
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			ProfileRegistration: received4,
			TestCase: TestCase{
				Response:   string(e4),
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for number, item := range cases {
		var buf bytes.Buffer
		var err error
		writer := multipart.NewWriter(&buf)

		if item.ProfileRegistration.Nickname != "" {
			err = writer.WriteField("nickname", item.ProfileRegistration.Nickname)
			if err != nil {
				t.Error(err)
			}
		}
		if item.ProfileRegistration.Email != "" {
			err = writer.WriteField("email", item.ProfileRegistration.Email)
			if err != nil {
				t.Error(err)
			}
		}
		if item.ProfileRegistration.Password != "" {
			err = writer.WriteField("password", item.ProfileRegistration.Password)
			if err != nil {
				t.Error(err)
			}
		}
		err = writer.Close()
		if err != nil {
			t.Error(err)
		}

		u := "http://127.0.0.1:8080/api/v1/profiles"

		req := httptest.NewRequest("POST", u, bytes.NewReader(buf.Bytes()))
		req.Header.Set("Content-Type", writer.FormDataContentType())

		w := httptest.NewRecorder()
		PostProfile(&env)(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				number, w.Code, item.StatusCode)
			return
		}

		resp := w.Result()
		var created models.Profile
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			t.Errorf("Process response error")
			return
		}

		if err := easyjson.UnmarshalFromReader(bytes.NewReader(body), &created); err == nil && resp.StatusCode == http.StatusOK {
			created.ID = 0
			body, err = easyjson.Marshal(created)
			if err != nil {
				t.Error(err)
			}
		}

		if string(body) != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				number, string(body), item.Response)
			return
		}
	}
}

func TestPutProfile(t *testing.T) {
	received1 := models.ProfileUpdate{
		ProfileInfo: models.ProfileInfo{
			Email:    "test2@mail.ru",
			Nickname: "test2",
			Password: "123456",
		},
		PasswordOld: "123",
	}

	received2 := models.ProfileUpdate{
		ProfileInfo: models.ProfileInfo{
			Email: "test3@mail.ru",
		},
	}

	received3 := models.ProfileUpdate{
		ProfileInfo: models.ProfileInfo{
			Nickname: "test3",
		},
	}

	r1, _ := easyjson.Marshal(received1)
	r2, _ := easyjson.Marshal(received2)
	r3, _ := easyjson.Marshal(received3)

	cases := []TestCasePutProfile{
		{
			ProfileID: 102,
			TestCase: TestCase{
				Body:       bytes.NewReader(r1),
				StatusCode: http.StatusNoContent,
			},
		},
		{
			ProfileID: 100,
			TestCase: TestCase{
				Body:       bytes.NewReader(r2),
				StatusCode: http.StatusBadRequest,
			},
		},
		{
			ProfileID: 100,
			TestCase: TestCase{
				Body:       bytes.NewReader(r3),
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for number, item := range cases {
		u := baseUrl + "/profiles/"
		vars := map[string]string{
			"id": strconv.FormatUint(item.ProfileID, 10),
		}
		req := httptest.NewRequest("PUT", u, item.Body)
		req = mux.SetURLVars(req, vars)
		ctx := context.WithValue(req.Context(), middleware.ProfileID, item.ProfileID)
		w := httptest.NewRecorder()
		PutProfile(&env)(w, req.WithContext(ctx))

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				number, w.Code, item.StatusCode)
			return
		}
	}
}

func TestGetProfiles(t *testing.T) {
	cases := []TestCaseGetProfiles{
		{
			Email:      "test3@mail.ru",
			StatusCode: http.StatusNoContent,
		},
		{
			Nickname:   "test3",
			StatusCode: http.StatusNoContent,
		},
		{
			Email:      "test4@mail.ru",
			StatusCode: http.StatusNotFound,
		},
		{
			Nickname:   "test4",
			StatusCode: http.StatusNotFound,
		},
		{
			StatusCode: http.StatusBadRequest,
		},
	}

	for number, item := range cases {
		u := baseUrl + "/profiles"
		if item.Email != "" {
			u = u + "?email=" + item.Email
		} else if item.Nickname != "" {
			u = u + "?nickname=" + item.Nickname
		}

		req := httptest.NewRequest("GET", u, nil)
		w := httptest.NewRecorder()
		GetProfiles(&env)(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				number, w.Code, item.StatusCode)
			return
		}
	}
}
