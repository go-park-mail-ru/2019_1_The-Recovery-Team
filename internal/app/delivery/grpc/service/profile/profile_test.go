package profile

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx"

	"github.com/magiconair/properties/assert"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/usecase"
)

var testCaseGet = []struct {
	name     string
	request  interface{}
	response interface{}
}{
	{
		name: "Test with not existing id",
		request: &GetRequest{
			Id: profile.NotExistingProfileId,
		},
		response: &GetResponse{},
	},
	{
		name: "Test with existing id",
		request: &GetRequest{
			Id: profile.ExistingProfileId,
		},
		response: &GetResponse{
			Info: &Info{
				Id:    profile.ExistingProfileId,
				Score: &Score{},
			},
		},
	},
}

func TestGet(t *testing.T) {
	service := Service{
		interactor: usecase.NewProfileInteractor(&profile.RepoMock{}),
	}
	for _, testCase := range testCaseGet {
		t.Run(testCase.name, func(t *testing.T) {
			response, _ := service.Get(context.Background(), testCase.request.(*GetRequest))
			assert.Equal(t, *response, *testCase.response.(*GetResponse), "Incorrect response")
		})
	}
}

var testCaseCreate = []struct {
	name     string
	request  interface{}
	response interface{}
}{
	{
		name: "Test with conflict data",
		request: &CreateRequest{
			Email:    profile.ExistingProfileEmail,
			Nickname: profile.ExistingProfileNickname,
			Password: profile.ExistingProfilePassword,
		},
		response: &CreateResponse{},
	},
	{
		name: "test with correct data",
		request: &CreateRequest{
			Email:    profile.NotExistingProfileEmail,
			Nickname: profile.NotExistingProfileNickname,
			Password: profile.NotExistingProfilePassword,
		},
		response: &CreateResponse{
			Id:       profile.DefaultProfileId,
			Email:    profile.CreatedProfileEmail,
			Nickname: profile.CreatedProfileNickname,
		},
	},
}

func TestCreate(t *testing.T) {
	service := Service{
		interactor: usecase.NewProfileInteractor(&profile.RepoMock{}),
	}
	for _, testCase := range testCaseCreate {
		t.Run(testCase.name, func(t *testing.T) {
			response, _ := service.Create(context.Background(), testCase.request.(*CreateRequest))
			assert.Equal(t, *response, *testCase.response.(*CreateResponse), "Incorrect response")
		})
	}
}

var testCaseUpdate = []struct {
	name     string
	request  interface{}
	response interface{}
}{
	{
		name: "Test with correct data",
		request: &UpdateRequest{
			Id:       profile.ExistingProfileId,
			Email:    profile.NotExistingProfileEmail,
			Nickname: profile.NotExistingProfileNickname,
		},
		response: &Nothing{},
	},
}

func TestUpdate(t *testing.T) {
	service := Service{
		interactor: usecase.NewProfileInteractor(&profile.RepoMock{}),
	}
	for _, testCase := range testCaseUpdate {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := service.Update(context.Background(), testCase.request.(*UpdateRequest))
			assert.Equal(t, *response, *testCase.response.(*Nothing), "Incorrect response")
			assert.Equal(t, err, nil, "Return error on correct data")
		})
	}
}

var testCaseUpdateAvatar = []struct {
	name     string
	request  interface{}
	response interface{}
}{
	{
		name: "Test with correct data",
		request: &UpdateAvatarRequest{
			Id:     profile.ExistingProfileId,
			Avatar: "",
		},
		response: &Nothing{},
	},
}

func TestUpdateAvatar(t *testing.T) {
	service := Service{
		interactor: usecase.NewProfileInteractor(&profile.RepoMock{}),
	}
	for _, testCase := range testCaseUpdateAvatar {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := service.UpdateAvatar(context.Background(), testCase.request.(*UpdateAvatarRequest))
			assert.Equal(t, *response, *testCase.response.(*Nothing), "Incorrect response")
			assert.Equal(t, err, nil, "Return error on correct data")
		})
	}
}

var testCaseGetByEmail = []struct {
	name     string
	request  interface{}
	response interface{}
	err      error
}{
	{
		name: "Test with forbidden email",
		request: &GetByEmailRequest{
			Email: profile.ForbiddenProfileEmail,
		},
		response: &GetResponse{},
		err:      errors.New(profile.DefaultError),
	},
	{
		name: "Test with not existing email",
		request: &GetByEmailRequest{
			Email: profile.NotExistingProfileEmail,
		},
		response: &GetResponse{},
		err:      pgx.ErrNoRows,
	},
	{
		name: "Test with correct data",
		request: &GetByEmailRequest{
			Email: profile.ExistingProfileEmail,
		},
		response: &GetResponse{
			Email: profile.ExistingProfileEmail,
			Info: &Info{
				Score: &Score{},
			},
		},
		err: nil,
	},
}

func TestGetByEmail(t *testing.T) {
	service := Service{
		interactor: usecase.NewProfileInteractor(&profile.RepoMock{}),
	}
	for _, testCase := range testCaseGetByEmail {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := service.GetByEmail(context.Background(), testCase.request.(*GetByEmailRequest))
			assert.Equal(t, *response, *testCase.response.(*GetResponse), "Incorrect response")
			assert.Equal(t, err, testCase.err, "Incorrect error value")
		})
	}
}

var testCaseGetByNickname = []struct {
	name     string
	request  interface{}
	response interface{}
	err      error
}{
	{
		name: "Test with forbidden nickname",
		request: &GetByNicknameRequest{
			Nickname: profile.ForbiddenProfileNickname,
		},
		response: &GetResponse{},
		err:      errors.New(profile.DefaultError),
	},
	{
		name: "Test with not existing nickname",
		request: &GetByNicknameRequest{
			Nickname: profile.NotExistingProfileNickname,
		},
		response: &GetResponse{},
		err:      pgx.ErrNoRows,
	},
	{
		name: "Test with correct data",
		request: &GetByNicknameRequest{
			Nickname: profile.ExistingProfileNickname,
		},
		response: &GetResponse{
			Info: &Info{
				Nickname: profile.ExistingProfileNickname,
				Score:    &Score{},
			},
		},
		err: nil,
	},
}

func TestGetByNickname(t *testing.T) {
	service := Service{
		interactor: usecase.NewProfileInteractor(&profile.RepoMock{}),
	}
	for _, testCase := range testCaseGetByNickname {
		t.Run(testCase.name, func(t *testing.T) {
			response, err := service.GetByNickname(context.Background(), testCase.request.(*GetByNicknameRequest))
			assert.Equal(t, *response, *testCase.response.(*GetResponse), "Incorrect response")
			assert.Equal(t, err, testCase.err, "Incorrect error value")
		})
	}
}

var testCaseGetByEmailAndPassword = []struct {
	name     string
	request  interface{}
	response interface{}
}{
	{
		name: "Test with not existing email",
		request: &GetByEmailAndPasswordRequest{
			Email:    profile.NotExistingProfileEmail,
			Password: profile.ExistingProfilePassword,
		},
		response: &GetResponse{},
	},
	{
		name: "Test with correct data",
		request: &GetByEmailAndPasswordRequest{
			Email:    profile.ExistingProfileEmail,
			Password: profile.ExistingProfilePassword,
		},
		response: &GetResponse{
			Info: &Info{
				Id:    profile.DefaultProfileId,
				Score: &Score{},
			},
		},
	},
}

func TestGetByEmailAndPassword(t *testing.T) {
	service := Service{
		interactor: usecase.NewProfileInteractor(&profile.RepoMock{}),
	}
	for _, testCase := range testCaseGetByEmailAndPassword {
		t.Run(testCase.name, func(t *testing.T) {
			response, _ := service.GetByEmailAndPassword(context.Background(), testCase.request.(*GetByEmailAndPasswordRequest))
			assert.Equal(t, *response, *testCase.response.(*GetResponse), "Incorrect response")
		})
	}
}

var testCaseList = []struct {
	name     string
	request  interface{}
	response interface{}
}{
	{
		name: "Test with incorrect data",
		request: &ListRequest{
			Limit: profile.ForbiddenLimit,
		},
		response: &ListResponse{},
	},
	{
		name: "Test with correct data",
		request: &ListRequest{
			Limit: 1,
		},
		response: &ListResponse{
			List: []*Info{
				{
					Id:    profile.DefaultProfileId,
					Score: &Score{},
				},
			},
		},
	},
}

func TestList(t *testing.T) {
	service := Service{
		interactor: usecase.NewProfileInteractor(&profile.RepoMock{}),
	}
	for _, testCase := range testCaseList {
		t.Run(testCase.name, func(t *testing.T) {
			response, _ := service.List(context.Background(), testCase.request.(*ListRequest))
			assert.Equal(t, *response, *testCase.response.(*ListResponse), "Incorrect response")
		})
	}
}

var testCaseCount = []struct {
	name     string
	request  interface{}
	response interface{}
}{
	{
		name:    "Test with correct data",
		request: &Nothing{},
		response: &CountResponse{
			Count: profile.DefaultCount,
		},
	},
}

func TestCount(t *testing.T) {
	service := Service{
		interactor: usecase.NewProfileInteractor(&profile.RepoMock{}),
	}
	for _, testCase := range testCaseCount {
		t.Run(testCase.name, func(t *testing.T) {
			response, _ := service.Count(context.Background(), testCase.request.(*Nothing))
			assert.Equal(t, *response, *testCase.response.(*CountResponse), "Incorrect response")
		})
	}
}
