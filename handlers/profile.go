package handlers

import (
	"api/database"
	"api/environment"
	"api/filesystem"
	"api/middleware"
	"github.com/asaskevich/govalidator"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"go.uber.org/zap"

	"api/models"

	"github.com/gorilla/mux"
)

func saveAvatar(env *environment.Env, avatar multipart.File, filename, dir string, id uint64) (string, error) {
	filename, err := filesystem.HashFileName(filename, id)
	if err != nil {
		return "", err
	}

	err = filesystem.SaveFile(avatar, dir, filename)
	if err != nil {
		return "", err
	}

	avatarPath := "/" + dir + filename
	err = env.Dbm.UpdateProfileAvatar(id, avatarPath)
	if err != nil {
		return "", err
	}
	return avatarPath, nil
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
func GetProfiles(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value("logger").(*zap.Logger)

		email := r.FormValue("email")
		nickname := r.FormValue("nickname")

		if email != "" && govalidator.IsEmail(email) {
			_, err := env.Dbm.GetProfileByEmail(email)
			if err != nil {
				if err == pgx.ErrNoRows {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				log.Error(err.Error(),
					zap.String("email", email),
				)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}

		if nickname != "" && govalidator.StringLength(nickname, "4", "20") {
			_, err := env.Dbm.GetProfileByNickname(nickname)
			if err != nil {
				if err == pgx.ErrNoRows {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				log.Error(err.Error(),
					zap.String("nickname", nickname),
				)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
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
func GetProfile(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value("logger").(*zap.Logger)

		vars := mux.Vars(r)
		id := vars["id"]

		profile, err := env.Dbm.GetProfile(id)
		if err != nil {
			if err == pgx.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			log.Error(err.Error(),
				zap.String("profile_id", id))
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

// PutProfile returns handler with environment which updates profile (email, nickname)
// @Summary Put profile
// @Description Update profile info
// @ID put-profile
// @Accept json
// @Param id path int true "Profile ID"
// @Param profile_info body models.ProfileUpdate true "Email, nickname"
// @Success 204 "Profile info is updated successfully"
// @Failure 400 "Incorrect request data"
// @Failure 403 "Not authorized"
// @Failure 404 "Not found"
// @Failure 422 {object} models.HandlerError "Invalid request data"
// @Failure 500 "Database error"
// @Router /profiles/{id} [PUT]
func PutProfile(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value("logger").(*zap.Logger)

		vars := mux.Vars(r)
		id, _ := strconv.ParseUint(vars["id"], 10, 64)

		profileID := r.Context().Value(middleware.ProfileID)
		if profileID != id {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		data := &models.ProfileUpdate{}
		err := unmarshalJSONBodyToStruct(r, data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if isValid, err := govalidator.ValidateStruct(data); !isValid && err != nil {
			message := models.HandlerError{
				Description: err.Error(),
			}
			writeResponseJSON(w, http.StatusUnprocessableEntity, message)
			return
		}

		if err = env.Dbm.UpdateProfile(id, data); err != nil {
			if err.Error() == "EmailAlreadyExists" || err.Error() == "NicknameAlreadyExists" {
				w.WriteHeader(http.StatusConflict)
				return
			}
			log.Error(err.Error(),
				zap.Uint64("profile_id", id),
				zap.String("email", data.Email),
				zap.String("nickname", data.Nickname))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// PutProfilePassword returns handler with environment which updates profile (email, nickname)
// @Summary Put profile password
// @Description Update profile password
// @ID put-profile_password
// @Accept json
// @Param id path int true "Profile ID"
// @Param profile_info body models.ProfileUpdatePassword true "Password"
// @Success 204 "Profile password is updated successfully"
// @Failure 400 "Incorrect request data"
// @Failure 403 "Not authorized"
// @Failure 404 "Not found"
// @Failure 422 {object} models.HandlerError "Invalid request data"
// @Failure 500 "Database error"
// @Router /profiles/{id}/password [PUT]
func PutProfilePassword(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value("logger").(*zap.Logger)

		vars := mux.Vars(r)
		id, _ := strconv.ParseUint(vars["id"], 10, 64)

		profileID := r.Context().Value(middleware.ProfileID)
		if profileID != id {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		data := &models.ProfileUpdatePassword{}
		err := unmarshalJSONBodyToStruct(r, data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if isValid, err := govalidator.ValidateStruct(data); !isValid && err != nil {
			message := models.HandlerError{
				Description: err.Error(),
			}
			writeResponseJSON(w, http.StatusUnprocessableEntity, message)
			return
		}

		data.Password, err = database.HashAndSalt(data.Password)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = env.Dbm.UpdateProfilePassword(id, data); err != nil {
			if err.Error() == "IncorrectProfilePassword" {
				message := models.HandlerError{
					Description: "Incorrect password",
				}
				writeResponseJSON(w, http.StatusUnprocessableEntity, message)
				return
			}
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
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
// @Param profile_info body models.ProfileCreate true "Email, nickname, password"
// @Param avatar body png false "Avatar"
// @Success 201 {object} models.ProfileCreated "Profile created successfully"
// @Failure 400 "Incorrect request data"
// @Failure 409 "Email or nickname already exists"
// @Failure 422 {object} models.HandlerError "Invalid request data"
// @Failure 500 "Database error"
// @Router /profiles [POST]
func PostProfile(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value("logger").(*zap.Logger)

		err := r.ParseMultipartForm(5 * (1 << 20)) // max size 5 MB
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		nickname := r.FormValue("nickname")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if nickname == "" || email == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data := &models.ProfileCreate{
			Email:    email,
			Nickname: nickname,
			Password: password,
		}

		if isValid, err := govalidator.ValidateStruct(data); !isValid && err != nil {
			message := models.HandlerError{
				Description: err.Error(),
			}
			writeResponseJSON(w, http.StatusUnprocessableEntity, message)
			return
		}

		data.Password, err = database.HashAndSalt(data.Password)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		created, err := env.Dbm.CreateProfile(data)
		if err != nil {
			if err.Error() == "EmailAlreadyExists" && err.Error() == "NicknameAlreadyExists" {
				log.Error(err.Error(),
					zap.String("email", data.Email),
					zap.String("nickname", data.Nickname))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			message := models.HandlerError{
				Description: err.Error(),
			}
			writeResponseJSON(w, http.StatusConflict, message)
			return
		}

		err = r.ParseMultipartForm(5 * (1 << 20)) // max size 5 MB
		if err != nil {
			writeResponseJSON(w, http.StatusCreated, created)
			return
		}

		avatar, header, err := r.FormFile("avatar")
		if err == nil {
			defer avatar.Close()
			filename := header.Filename
			dir := "upload/img/"

			if avatarPath, err := saveAvatar(env, avatar, filename, dir, created.ID); err == nil {
				created.Avatar = avatarPath
			}
		}

		token, err := env.Sm.Set(created.ID, 24*time.Hour)
		if err != nil {
			log.Error(err.Error(),
				zap.Uint64("created_id", created.ID))
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

		writeResponseJSON(w, http.StatusCreated, created)
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
func PutAvatar(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value("logger").(*zap.Logger)
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
			log.Error(err.Error(),
				zap.Uint64("profile_id", id))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := &models.ProfileAvatar{
			Avatar: avatarPath,
		}
		writeResponseJSON(w, http.StatusOK, result)
	}
}
