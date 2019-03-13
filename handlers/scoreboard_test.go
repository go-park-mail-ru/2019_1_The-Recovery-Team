package handlers

import (
	"api/models"
	"github.com/mailru/easyjson"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type TestCaseGetScoreboard struct {
	Limit  int
	Offset int
	TestCase
}

func TestGetScoreboard(t *testing.T) {
	var total int
	err := env.Dbm.Find(&total, QueryCountProfilesNumber)
	if err != nil {
		t.Error("Postgres connection error")
		return
	}

	expected1 := models.Profiles{
		List:  []models.Profile{},
		Total: total,
	}
	err = env.Dbm.FindAll(&expected1.List, QueryProfilesWithLimitAndOffset, 100, 1)
	if err != nil {
		t.Error("Postgres connection error")
		return
	}

	expected2 := models.Profiles{
		List:  []models.Profile{},
		Total: total,
	}
	err = env.Dbm.FindAll(&expected2.List, QueryProfilesWithOffset, 1)
	if err != nil {
		t.Error("Postgres connection error")
		return
	}

	expected3 := models.Profiles{
		List:  []models.Profile{},
		Total: total,
	}
	err = env.Dbm.FindAll(&expected3.List, QueryProfilesWithLimit, 1)
	if err != nil {
		t.Error("Postgres connection error")
		return
	}

	e1, _ := easyjson.Marshal(expected1)
	e2, _ := easyjson.Marshal(expected2)
	e3, _ := easyjson.Marshal(expected3)

	cases := []TestCaseGetScoreboard{
		{
			Limit:  100,
			Offset: 1,
			TestCase: TestCase{
				Response:   string(e1),
				StatusCode: http.StatusOK,
			},
		},
		{
			Offset: 1,
			TestCase: TestCase{
				Response:   string(e2),
				StatusCode: http.StatusOK,
			},
		},
		{
			Limit: 1,
			TestCase: TestCase{
				Response:   string(e3),
				StatusCode: http.StatusOK,
			},
		},
		{
			TestCase: TestCase{
				Response:   "",
				StatusCode: http.StatusBadRequest,
			},
		},
	}

	for number, item := range cases {
		u := "http://127.0.0.1:8080/api/v1/scores"
		if item.Limit != 0 && item.Offset != 0 {
			u = u + "?start=" + strconv.Itoa(item.Offset) + "&limit=" + strconv.Itoa(item.Limit)
		} else if item.Limit != 0 {
			u = u + "?limit=" + strconv.Itoa(item.Limit)
		} else if item.Offset != 0 {
			u = u + "?start=" + strconv.Itoa(item.Offset)
		}
		req := httptest.NewRequest("GET", u, nil)
		w := httptest.NewRecorder()
		GetScoreboard(&env)(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong StatusCode: got %d, expected %d",
				number, w.Code, item.StatusCode)
			return
		}

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
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
