package handlers

import (
	"api/middleware"
	"api/models"
	"bytes"
	"context"
	"github.com/mailru/easyjson"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestCaseGetSession struct {
	ProfileID uint64
	TestCase
}

type TestCaseDeleteSession struct {
	SessionID string
	TestCase
}

func TestGetSession(t *testing.T) {
	_, err := env.Sm.Set(100, 24*time.Hour)
	if err != nil {
		t.Error("Redis connection error")
	}

	expected := models.Profile{
		ID: 100,
	}

	e, _ := easyjson.Marshal(expected)

	cases := []TestCaseGetSession{
		{
			ProfileID: 100,
			TestCase: TestCase{
				Response:   string(e),
				StatusCode: http.StatusOK,
			},
		},
	}

	for number, item := range cases {
		u := "http://127.0.0.1:8080/api/v1/sessions"
		req := httptest.NewRequest("GET", u, nil)
		ctx := context.WithValue(req.Context(), middleware.ProfileID, item.ProfileID)
		w := httptest.NewRecorder()

		GetSession(&env)(w, req.WithContext(ctx))

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				number, w.Code, item.StatusCode)
			return
		}

		resp := w.Result()
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Process response error")
			return
		}
		if string(body) != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				number, string(body), item.Response)
			return
		}
	}
}

func TestPostSession(t *testing.T) {
	expected := models.Profile{
		ID: 100,
		ProfileInfo: models.ProfileInfo{
			Email:    "test@mail.ru",
			Nickname: "test",
		},
	}
	received := models.ProfileLogin{
		Email:    "test@mail.ru",
		Password: "123",
	}

	r, _ := easyjson.Marshal(received)
	e, _ := easyjson.Marshal(expected)

	cases := []TestCase{
		{
			Body:       bytes.NewReader(r),
			Response:   string(e),
			StatusCode: http.StatusOK,
		},
	}

	for number, item := range cases {
		u := "http://127.0.0.1:8080/api/v1/sessions"
		req := httptest.NewRequest("POST", u, item.Body)
		w := httptest.NewRecorder()

		PostSession(&env)(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				number, w.Code, item.StatusCode)
			return
		}

		resp := w.Result()
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Process response error")
			return
		}
		if string(body) != item.Response {
			t.Errorf("[%d] wrong Response: got %+v, expected %+v",
				number, string(body), item.Response)
			return
		}
	}
}

func TestDeleteSession(t *testing.T) {
	sessionID, err := env.Sm.Set(100, 24*time.Hour)
	if err != nil {
		t.Error("Redis connection error")
	}

	cases := []TestCaseDeleteSession{
		{
			SessionID: sessionID,
			TestCase: TestCase{
				StatusCode: http.StatusOK,
			},
		},
	}

	for number, item := range cases {
		u := "http://127.0.0.1:8080/api/v1/sessions"
		req := httptest.NewRequest("DELETE", u, nil)
		ctx := context.WithValue(req.Context(), middleware.SessionID, item.SessionID)
		w := httptest.NewRecorder()

		req.AddCookie(&http.Cookie{
			Name:     "session_id",
			Value:    item.SessionID,
			Path:     "/",
			Expires:  time.Now().Add(24*time.Hour - 10*time.Minute),
			HttpOnly: true,
		})

		DeleteSession(&env)(w, req.WithContext(ctx))

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				number, w.Code, item.StatusCode)
			return
		}
	}
}
