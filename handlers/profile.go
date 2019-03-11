package handlers

import (
	"api/filesystem"
	"api/middleware"
	"database/sql"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

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

func updateProfile(env *models.Env, id uint64, newInfo *models.ProfileUpdate) error {
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
		passwordHash, err := hashAndSalt(newInfo.Password)
		if err != nil {
			return err
		}
		newInfo.Password = passwordHash
		if set != "" {
			set = set + ", "
		}
		set = set + "password = :password"
	}

	dbo := env.Dbm.DB()
	if set != "" {
		query := `UPDATE profile SET ` + set + ` WHERE id = :id`
		_, err := dbo.NamedExec(query, &struct {
			ID uint64
			models.ProfileUpdate
		}{
			ID:            id,
			ProfileUpdate: *newInfo,
		})
		return err
	}
	return nil
}

func checkField(env *models.Env, w http.ResponseWriter, table, field, value string) {
	exists, _ := env.Dbm.FindWithField(table, field, value)
	if exists {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

// GetProfiles returns handler with environment which processes request for checking email or nickname existens
// @Summary Get profiles
// @Description Check profile existence with email or nickname
// @ID get-profiles
// @Param email path string false "Profile email"
// @Param nickname path string false "Profile nickname"
// @Success 204 "Profile found successfully"
// @Failure 400 "Incorrect request data"
// @Failure 403 "Not authorized"
// @Failure 404 "Not found"
// @Failure 500 "Database error"
// @Router /profiles [GET]
func GetProfiles(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		nickname := r.FormValue("nickname")

		if email != "" {
			checkField(env, w, "profile", "email", email)
			return
		}

		if nickname != "" {
			checkField(env, w, "profile", "nickname", nickname)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}

// GetProfile returns handler with environment which processes request for getting profile by id
// @Summary Get profile
// @Description Get profile info (for profile owner returns info with email)
// @ID get-profile
// @Produce json
// @Param id path int true "Profile ID"
// @Success 200 {object} models.Profile "Profile found successfully"
// @Failure 403 "Not authorized"
// @Failure 404 "Not found"
// @Failure 500 "Database error"
// @Router /profiles/{id} [GET]
func GetProfile(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseUint(vars["id"], 10, 64)

		profile := &models.Profile{}
		err := env.Dbm.Find(profile, QueryProfileById, id)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return

			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		profileID := r.Context().Value(middleware.ProfileID)
		if profileID != id {
			profile.Email = ""
		}

		writeResponseJSON(w, http.StatusOK, profile)
	}
}

// PutProfile returns handler with environment which updates profile
// @Summary Put profile
// @Description Update profile info
// @ID put-profile
// @Accept json
// @Param id path int true "Profile ID"
// @Param profile_info body models.ProfileInfo true "Email, nickname, password"
// @Success 204 "Profile info is updated successfully"
// @Failure 400 "Incorrect request data"
// @Failure 403 "Not authorized"
// @Failure 404 "Not found"
// @Failure 422 "Incorrrect current password"
// @Failure 500 "Database error"
// @Router /profiles/{id} [PUT]
func PutProfile(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, _ := strconv.ParseUint(vars["id"], 10, 64)

		profileID := r.Context().Value(middleware.ProfileID)
		if profileID != id {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		newInfo := &models.ProfileUpdate{}
		err := unmarshalJSONBodyToStruct(r, newInfo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newInfo.Password != "" {
			if newInfo.PasswordOld != "" {
				var currentPassword string
				err := env.Dbm.Find(&currentPassword, QueryProfilePassword, id)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if matches, err := verifyPassword(newInfo.PasswordOld, currentPassword); !matches || err != nil {
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

		}

		err = updateProfile(env, id, newInfo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// PostProfile returns handler with environment which creates profile
// @Summary Post profile
// @Description Create profile
// @ID post-profile
// @Accept multipart/form-data
// @Produce json
// @Param profile_info body models.ProfileInfo true "Email, nickname, password"
// @Param avatar body png false "Avatar"
// @Success 200 {object} models.Profile "Profile created successfully"
// @Failure 400 "Incorrect request data"
// @Failure 500 "Database error"
// @Router /profiles [POST]
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
			Nickname: r.FormValue("nickname"),
			ProfileLogin: models.ProfileLogin{
				Email:    r.FormValue("email"),
				Password: r.FormValue("password"),
			},
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
		if err == nil {
			defer avatar.Close()
			filename := header.Filename
			dir := "upload/img/"

			if avatarPath, err := saveAvatar(env, avatar, filename, dir, result.ID); err == nil {
				result.Avatar = avatarPath
			}
		}

		token, err := env.Sm.Set(result.ID, 24*time.Hour)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(24*time.Hour - 10*time.Minute),
			HttpOnly: true,
		})

		writeResponseJSON(w, http.StatusOK, result)
	}
}

// PutAvatar returns handler with environment which adds or updates profile avatar
// @Summary Put avatar
// @Description Update profile avatar
// @ID put-avatar
// @Accept multipart/form-data
// @Produce json
// @Param avatar body png false "Avatar"
// @Success 200 {object} models.ProfileAvatar "Profile avatar is updated successfully"
// @Failure 400 "Incorrect request data"
// @Failure 403 "Not authorized"
// @Failure 500 "Database error"
// @Router /avatars [PUT]
func PutAvatar(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(middleware.ProfileID).(uint64)

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

		avatarPath, err := saveAvatar(env, avatar, filename, dir, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := &models.ProfileAvatar{
			Avatar: avatarPath,
		}
		writeResponseJSON(w, http.StatusOK, result)
	}
}
