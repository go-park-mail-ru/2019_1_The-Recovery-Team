package handlers

import (
	"api/filesystem"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/mailru/easyjson"
	"golang.org/x/crypto/bcrypt"

	"api/models"

	"github.com/gorilla/mux"
)

func saveAvatar(env *models.Env, avatar multipart.File, filename, dir string, id uint64) (string, error) {
	filename, err := filesystem.HashFileName(filename, id)
	if err != nil {
		return "", err
	}

	err = filesystem.SaveFile(avatar, dir, filename)
	if err != nil {
		return "", err
	}

	avatarPath := "/" + dir + filename
	err = env.Dbm.Update(QueryUpdateProfileAvatar, avatarPath, id)
	if err != nil {
		return "", err
	}
	return avatarPath, nil
}

func hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func updateProfile(env *models.Env, id uint64, newInfo *models.ProfileInfo) error {
	var set string
	if newInfo.Nickname != "" {
		exists, _ := env.Dbm.FindWithField("profile", "nickname", newInfo.Nickname)
		if exists {
			return errors.New("nickname already exists")
		}
		set = set + "nickname = :nickname"
	}
	if newInfo.Email != "" {
		exists, _ := env.Dbm.FindWithField("profile", "email", newInfo.Email)
		if exists {
			return errors.New("email already exists")
		}
		if set != "" {
			set = set + ", "
		}
		set = set + "email = :email"
	}
	if newInfo.Password != "" {
		if set != "" {
			set = set + ", "
		}
		set = set + "password = :password"
	}

	dbo := env.Dbm.DB()
	if set != "" {
		query := `UPDATE profile SET ` + set + ` WHERE id = :id`
		_, err := dbo.NamedExec(query, &models.Profile{
			ID:          id,
			ProfileInfo: *newInfo,
		})
		return err
	}
	return nil
}

func checkFieldHandler(env *models.Env, table, field string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		value, _ := vars[field]
		exists, _ := env.Dbm.FindWithField(table, field, value)
		if exists {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

// GetProfile returns handler with environment which processes request for getting profile by id
func GetProfile(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := vars["id"]

		profile := &models.Profile{}
		err := env.Dbm.Find(profile, QueryProfileById, id)
		if err != nil {
			message := models.HandlerError{
				Description: "user doesn't exist",
			}
			writeResponseJSON(w, http.StatusNotFound, message)
			return
		}

		writeResponseJSON(w, http.StatusOK, profile)
	}
}

// CheckProfileEmail returns handler with environment which checks existence of profile email
func CheckProfileEmail(env *models.Env) http.HandlerFunc {
	return checkFieldHandler(env, "profile", "email")
}

// CheckProfileNickname returns handler with environment which checks existence of profile nickname
func CheckProfileNickname(env *models.Env) http.HandlerFunc {
	return checkFieldHandler(env, "profile", "nickname")
}

// GetProfiles returns handler with environment which processes request for getting profiles order by score
func GetProfiles(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, limitErr := strconv.ParseInt(r.FormValue("limit"), 10, 64)
		offset, offsetErr := strconv.ParseInt(r.FormValue("start"), 10, 64)

		result := models.Profiles{}
		var err error
		switch {
		case limitErr == nil && offsetErr == nil:
			{
				env.Dbm.FindAll(&result, QueryProfilesWithLimitAndOffset, limit, offset)
			}
		case limitErr == nil:
			{
				env.Dbm.FindAll(&result, QueryProfilesWithLimit, limit)
			}
		case offsetErr == nil:
			{
				env.Dbm.FindAll(&result, QueryProfilesWithOffset, offset)
			}
		default:
			{
				env.Dbm.FindAll(&result, QueryProfiles)
			}
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			message := models.HandlerError{
				Description: "server error",
			}
			easyjson.MarshalToHTTPResponseWriter(message, w)
			return
		}

		w.WriteHeader(http.StatusOK)
		easyjson.MarshalToHTTPResponseWriter(result, w)
	}
}

// PutProfile returns handler with environment which updates profile
func PutProfile(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseUint(vars["id"], 10, 64)

		newInfo := &models.ProfileInfo{}
		err := unmarshalJSONBodyToStruct(r, newInfo)
		if err != nil {
			message := models.HandlerError{
				Description: "incorrect data",
			}
			writeResponseJSON(w, http.StatusBadRequest, message)
			return
		}

		err = updateProfile(env, id, newInfo)
		if err != nil {
			message := models.HandlerError{
				Description: err.Error(),
			}
			writeResponseJSON(w, http.StatusForbidden, message)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// PostProfile returns handler with environment which creates profile
func PostProfile(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(5 * (1 << 20)) // max size 5 MB
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if r.FormValue("nickname") == "" || r.FormValue("email") == "" || r.FormValue("password") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newProfile := &models.ProfileRegistration{
			Email:    r.FormValue("email"),
			Nickname: r.FormValue("nickname"),
			Password: r.FormValue("password"),
		}

		if exists, err := env.Dbm.FindWithField("profile", "email", newProfile.Email); err != nil || exists {
			message := models.HandlerError{
				Description: "email already exists",
			}
			writeResponseJSON(w, http.StatusBadRequest, message)
			return
		}

		if exists, err := env.Dbm.FindWithField("profile", "nickname", newProfile.Nickname); err != nil || exists {
			message := models.HandlerError{
				Description: "nickname already exists",
			}
			writeResponseJSON(w, http.StatusBadRequest, message)
			return
		}

		password, err := hashAndSalt(newProfile.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := &models.Profile{}
		err = env.Dbm.Create(result, QueryInsertProfile, newProfile.Email, newProfile.Nickname, password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = r.ParseMultipartForm(5 * (1 << 20)) // max size 5 MB
		if err != nil {
			writeResponseJSON(w, http.StatusOK, result)
			return
		}

		avatar, header, err := r.FormFile("avatar")
		defer avatar.Close()
		if err == nil {
			filename := header.Filename
			dir := "upload/img/"

			if avatarPath, err := saveAvatar(env, avatar, filename, dir, result.ID); err == nil {
				result.Avatar = avatarPath
			}
		}

		writeResponseJSON(w, http.StatusOK, result)
	}
}

// PutAvatar adds or updates profile avatar
func PutAvatar(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseUint(vars["id"], 10, 64)

		err := r.ParseMultipartForm(5 * (1 << 20)) // max size 5 MB
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		avatar, header, err := r.FormFile("avatar")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer avatar.Close()
		filename := header.Filename
		dir := "upload/img/"

		if _, err := saveAvatar(env, avatar, filename, dir, id); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
