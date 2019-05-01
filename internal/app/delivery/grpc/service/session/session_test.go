package session

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/redis/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/stretchr/testify/assert"
)

var testCaseGet = []struct {
	name     string
	request  interface{}
	response interface{}
}{
	{
		name: "Test with not existing session id",
		request: &SessionId{
			Id: session.Unauthorized,
		},
		response: &ProfileId{},
	},
	{
		name: "Test with correct data",
		request: &SessionId{
			Id: session.Authorized,
		},
		response: &ProfileId{
			Id: session.AuthorizedProfileId,
		},
	},
}

func TestGet(t *testing.T) {
	service := Service{
		interactor: usecase.NewSessionInteractor(&session.RepoMock{}),
	}
	for _, testCase := range testCaseGet {
		t.Run(testCase.name, func(t *testing.T) {
			response, _ := service.Get(context.Background(), testCase.request.(*SessionId))
			assert.Equal(t, *response, *testCase.response.(*ProfileId), "Incorrect response")
		})
	}
}

var testCaseSet = []struct {
	name     string
	request  interface{}
	response interface{}
}{
	{
		name: "Test with correct data",
		request: &Create{
			ProfileId: &ProfileId{
				Id: session.DefaultProfileId,
			},
			Expires: &duration.Duration{},
		},
		response: &SessionId{
			Id: session.Authorized,
		},
	},
}

func TestSet(t *testing.T) {
	service := Service{
		interactor: usecase.NewSessionInteractor(&session.RepoMock{}),
	}
	for _, testCase := range testCaseSet {
		t.Run(testCase.name, func(t *testing.T) {
			response, _ := service.Set(context.Background(), testCase.request.(*Create))
			assert.Equal(t, *response, *testCase.response.(*SessionId), "Incorrect response")
		})
	}
}

var testCaseDelete = []struct {
	name     string
	request  interface{}
	response interface{}
}{
	{
		name: "Test with correct data",
		request: &SessionId{
			Id: session.Authorized,
		},
		response: &Nothing{},
	},
}

func TestDelete(t *testing.T) {
	service := Service{
		interactor: usecase.NewSessionInteractor(&session.RepoMock{}),
	}
	for _, testCase := range testCaseDelete {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := service.Delete(context.Background(), testCase.request.(*SessionId))
			assert.Equal(t, *response, *testCase.response.(*Nothing), "Incorrect response")
			assert.Equal(t, err, nil, "Incorrect error value")
		})
	}
}
