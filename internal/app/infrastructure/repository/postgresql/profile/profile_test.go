package profile

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/profile"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql"

	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
)

func repo() Repo {
	conn := postgresql.ConnMock{
		Tx: postgresql.TxMock{},
	}
	return Repo{
		conn: &conn,
	}
}

func TestNewProfileRepo(t *testing.T) {
	conn := &pgx.ConnPool{}
	assert.NotEmpty(t, NewRepo(conn),
		"Doesn't create profile repository instance")
}

var testCaseGet = []struct {
	name string
	id   uint64
	err  error
}{
	{
		name: "Test with not existing id",
		id:   postgresql.ForbiddenProfileId,
		err:  errors.New(postgresql.DefaultError),
	},
	{
		name: "Test with existing id",
		id:   postgresql.ProfileId,
		err:  nil,
	},
}

func TestGet(t *testing.T) {
	repo := repo()

	for _, testCase := range testCaseGet {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.Get(testCase.id)
			assert.Equal(t, testCase.err, err, "Return incorrect error value")
		})
	}
}

var testCaseCreate = []struct {
	name   string
	create profile.Create
	err    error
}{
	{
		name: "Test with conflict email",
		create: profile.Create{
			Email:    postgresql.ConflictProfileEmail,
			Nickname: postgresql.ProfileNickname,
			Password: postgresql.ProfilePassword,
		},
		err: errors.New(EmailAlreadyExists),
	},
	{
		name: "Test with conflict nickname",
		create: profile.Create{
			Email:    postgresql.ProfileEmail,
			Nickname: postgresql.ConflictProfileNickname,
			Password: postgresql.ProfilePassword,
		},
		err: errors.New(NicknameAlreadyExists),
	},
	{
		name: "Test with correct data",
		create: profile.Create{
			Email:    postgresql.ProfileEmail,
			Nickname: postgresql.ProfileNickname,
			Password: postgresql.ProfilePassword,
		},
		err: nil,
	},
}

func TestCreate(t *testing.T) {
	repo := repo()

	for _, testCase := range testCaseCreate {
		t.Run(testCase.name, func(t *testing.T) {
			prof, err := repo.Create(&testCase.create)
			assert.Equal(t, testCase.err, err, "Return incorrect error value")
			fmt.Println(prof)
		})
	}
}

var testCaseUpdate = []struct {
	name   string
	id     uint64
	update profile.UpdateInfo
	err    error
}{
	{
		name: "Test with conflict email",
		id:   postgresql.ProfileId,
		update: profile.UpdateInfo{
			Email:    postgresql.ConflictProfileEmail,
			Nickname: postgresql.ProfileNickname,
		},
		err: errors.New(EmailAlreadyExists),
	},
	{
		name: "Test with conflict nickname",
		id:   postgresql.ProfileId,
		update: profile.UpdateInfo{
			Email:    postgresql.ProfileEmail,
			Nickname: postgresql.ConflictProfileNickname,
		},
		err: errors.New(NicknameAlreadyExists),
	},
	{
		name: "Test with not existing id",
		id:   postgresql.ForbiddenProfileId,
		update: profile.UpdateInfo{
			Email:    postgresql.ProfileEmail,
			Nickname: postgresql.ProfileNickname,
		},
		err: errors.New(DefaultError),
	},
	{
		name: "Test with correct data",
		id:   postgresql.ProfileId,
		update: profile.UpdateInfo{
			Email:    postgresql.ProfileEmail,
			Nickname: postgresql.ProfileNickname,
		},
		err: nil,
	},
}

func TestUpdate(t *testing.T) {
	repo := repo()

	for _, testCase := range testCaseUpdate {
		t.Run(testCase.name, func(t *testing.T) {
			err := repo.Update(testCase.id, &testCase.update)
			assert.Equal(t, testCase.err, err, "Return incorrect error value")
		})
	}
}

func TestUpdateAvatar(t *testing.T) {
	repo := repo()
	err := repo.UpdateAvatar(postgresql.ProfileId, postgresql.ProfileAvatar)
	assert.Empty(t, err, "Return error on correct data")
}

var testCaseUpdatePassword = []struct {
	name   string
	id     uint64
	update profile.UpdatePassword
	err    error
}{
	{
		name: "Test with not existing user",
		id:   postgresql.ForbiddenProfileId,
		update: profile.UpdatePassword{
			Password:    postgresql.ProfilePassword,
			PasswordOld: postgresql.ProfilePassword,
		},
		err: errors.New(DefaultError),
	},
	{
		name: "Test with incorrect old password",
		id:   postgresql.ProfileId,
		update: profile.UpdatePassword{
			Password:    postgresql.ProfilePassword,
			PasswordOld: "",
		},
		err: errors.New(IncorrectProfilePassword),
	},
	{
		name: "Test with correct data",
		id:   postgresql.ProfileId,
		update: profile.UpdatePassword{
			Password:    postgresql.ProfilePassword,
			PasswordOld: postgresql.ProfilePassword,
		},
		err: nil,
	},
}

func TestUpdatePassword(t *testing.T) {
	repo := repo()

	for _, testCase := range testCaseUpdatePassword {
		t.Run(testCase.name, func(t *testing.T) {
			err := repo.UpdatePassword(testCase.id, &testCase.update)
			assert.Equal(t, testCase.err, err, "Return incorrect error value")
		})
	}
}

func TestGetByEmail(t *testing.T) {
	repo := repo()
	_, err := repo.GetByEmail(postgresql.ProfileEmail)
	assert.Empty(t, err, "Return error on correct data")
}

func TestGetByNickname(t *testing.T) {
	repo := repo()
	_, err := repo.GetByNickname(postgresql.ProfileNickname)
	assert.Empty(t, err, "Return error on correct data")
}

var testCaseGetByEmailAndPassword = []struct {
	name string
	data profile.Login
	err  error
}{
	{
		name: "Test with not existing user",
		data: profile.Login{
			Email:    postgresql.ForbiddenEmail,
			Password: postgresql.ProfilePassword,
		},
		err: errors.New(DefaultError),
	},
	{
		name: "Test with incorrect password",
		data: profile.Login{
			Email:    postgresql.ProfileEmail,
			Password: "",
		},
		err: pgx.ErrNoRows,
	},
	{
		name: "Test with correct data",
		data: profile.Login{
			Email:    postgresql.ProfileEmail,
			Password: postgresql.ProfilePassword,
		},
		err: nil,
	},
}

func TestGetByEmailAndPassword(t *testing.T) {
	repo := repo()

	for _, testCase := range testCaseGetByEmailAndPassword {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.GetByEmailAndPassword(&testCase.data)
			assert.Equal(t, testCase.err, err, "Return incorrect error value")
		})
	}
}

var testCaseList = []struct {
	name   string
	limit  int64
	offset int64
	err    error
}{
	{
		name:  "Test with incorrect limit",
		limit: postgresql.ForbiddenLimit,
		err:   errors.New(DefaultError),
	},
	{
		name:  "Test with correct data",
		limit: 1,
		err:   nil,
	},
}

func TestList(t *testing.T) {
	repo := repo()

	for _, testCase := range testCaseList {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.List(testCase.limit, testCase.offset)
			assert.Equal(t, testCase.err, err, "Return incorrect error value")
		})
	}
}

func TestCount(t *testing.T) {
	repo := repo()

	_, err := repo.Count()
	assert.Empty(t, err, "Return error on correct data")
}
