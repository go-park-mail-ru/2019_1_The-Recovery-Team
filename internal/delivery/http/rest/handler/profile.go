package handler

import (
	"mime/multipart"
	"net/http"
	handler "sadislands/internal/delivery/http/rest/handler/error"
	"sadislands/internal/delivery/http/rest/handler/saver"
	"sadislands/internal/delivery/http/rest/handler/unmarshaler"
	"sadislands/internal/delivery/http/rest/handler/writer"
	"sadislands/internal/delivery/http/rest/middleware"
	"sadislands/internal/domain/profile"
	"sadislands/internal/infrastructure/repository/postgresql"
	"sadislands/internal/usecase"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgx"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const (
	EmailAlreadyExists       = "EmailAlreadyExists"
	NicknameAlreadyExists    = "NicknameAlreadyExists"
	IncorrectProfilePassword = "IncorrectProfilePassword"
)

func saveAvatar(profileInteractor *usecase.ProfileInteractor, avatar multipart.File, filename, dir string, id uint64) (string, error) {
	// Hash filename
	filename, err := saver.HashFileName(filename, id)
	if err != nil {
		return "", err
	}

	err = saver.SaveFile(avatar, dir, filename)
	if err != nil {
		return "", err
	}

	// Updates profile avatar path in database
	avatarPath := "/" + dir + filename
	err = profileInteractor.UpdateProfileAvatar(id, avatarPath)
	if err != nil {
		return "", err
	}
	return avatarPath, nil
}

// GetProfiles returns handler with environment which processes request for checking email or nickname existence
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
func GetProfiles(profileInteractor *usecase.ProfileInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)

		email := r.FormValue("email")
		nickname := r.FormValue("nickname")

		// Check email existence
		if email != "" && govalidator.IsEmail(email) {
			_, err := profileInteractor.GetProfileByEmail(email)
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

		// Check nickname existence
		if nickname != "" && govalidator.StringLength(nickname, "4", "20") {
			_, err := profileInteractor.GetProfileByNickname(nickname)
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
// @Success 200 {object} profile.Profile "Profile found successfully"
// @Failure 403 "Not authorized"
// @Failure 404 "Not found"
// @Failure 500 "Database error"
// @Router /profiles/{id} [GET]
func GetProfile(profileInteractor *usecase.ProfileInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)

		id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		prof, err := profileInteractor.GetProfile(id)
		if err != nil {
			if err == pgx.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			log.Error(err.Error(),
				zap.Uint64("profile_id", id))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		profileID := r.Context().Value(middleware.ProfileID)
		if profileID != id {
			prof.Email = ""
		}

		writer.WriteResponseJSON(w, http.StatusOK, prof)
	}
}

// PutProfile returns handler with environment which updates profile (email, nickname)
// @Summary Put profile
// @Description Update profile info
// @ID put-profile
// @Accept json
// @Param id path int true "Profile ID"
// @Param profile_info body profile.UpdateInfo true "Email, nickname"
// @Success 204 "Profile info is updated successfully"
// @Failure 400 "Incorrect request data"
// @Failure 403 "Not authorized"
// @Failure 404 "Not found"
// @Failure 422 {object} handler.Error "Invalid request data"
// @Failure 500 "Database error"
// @Router /profiles/{id} [PUT]
func PutProfile(profileInteractor *usecase.ProfileInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)

		id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Verification of rights for this profile
		profileID := r.Context().Value(middleware.ProfileID)
		if profileID != id {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		data := &profile.UpdateInfo{}
		err = unmarshaler.UnmarshalJSONBodyToStruct(r, data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if isValid, err := govalidator.ValidateStruct(data); !isValid && err != nil {
			message := handler.Error{
				Description: err.Error(),
			}
			writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, message)
			return
		}

		if err = profileInteractor.UpdateProfile(id, data); err != nil {
			if err.Error() == EmailAlreadyExists || err.Error() == NicknameAlreadyExists {
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
// @Param profile_info body profile.UpdatePassword true "Password"
// @Success 204 "Profile password is updated successfully"
// @Failure 400 "Incorrect request data"
// @Failure 403 "Not authorized"
// @Failure 404 "Not found"
// @Failure 422 {object} handler.Error "Invalid request data"
// @Failure 500 "Database error"
// @Router /profiles/{id}/password [PUT]
func PutProfilePassword(profileInteractor *usecase.ProfileInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)

		id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Verification of rights for this profile
		profileID := r.Context().Value(middleware.ProfileID)
		if profileID != id {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		data := &profile.UpdatePassword{}
		err = unmarshaler.UnmarshalJSONBodyToStruct(r, data)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if isValid, err := govalidator.ValidateStruct(data); !isValid && err != nil {
			message := handler.Error{
				Description: err.Error(),
			}
			writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, message)
			return
		}

		data.Password, err = postgresql.HashAndSalt(data.Password)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = profileInteractor.UpdateProfilePassword(id, data); err != nil {
			if err.Error() == IncorrectProfilePassword {
				message := handler.Error{
					Description: "Incorrect password",
				}
				writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, message)
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
// @Param profile_info body profile.Create true "Email, nickname, password"
// @Param avatar body png false "Avatar"
// @Success 201 {object} profile.Created "Profile created successfully"
// @Failure 400 "Incorrect request data"
// @Failure 409 "Email or nickname already exists"
// @Failure 422 {object} handler.Error "Invalid request data"
// @Failure 500 "Database error"
// @Router /profiles [POST]
func PostProfile(profileInteractor *usecase.ProfileInteractor, sessionInteractor *usecase.SessionInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)

		err := r.ParseMultipartForm(5 * (1 << 20)) // Max avatar size 5 MB
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

		data := &profile.Create{
			Email:    email,
			Nickname: nickname,
			Password: password,
		}

		if isValid, err := govalidator.ValidateStruct(data); !isValid && err != nil {
			message := handler.Error{
				Description: err.Error(),
			}
			writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, message)
			return
		}

		data.Password, err = postgresql.HashAndSalt(data.Password)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create profile
		created, err := profileInteractor.CreateProfile(data)
		if err != nil {
			if err.Error() == EmailAlreadyExists && err.Error() == NicknameAlreadyExists {
				log.Error(err.Error(),
					zap.String("email", data.Email),
					zap.String("nickname", data.Nickname))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			message := handler.Error{
				Description: err.Error(),
			}
			writer.WriteResponseJSON(w, http.StatusConflict, message)
			return
		}

		// Save profile avatar
		avatar, header, err := r.FormFile("avatar")
		if err == nil {
			defer avatar.Close()
			filename := header.Filename
			dir := "upload/img/"

			if avatarPath, err := saveAvatar(profileInteractor, avatar, filename, dir, created.ID); err == nil {
				created.Avatar = avatarPath
			} else {
				log.Warn(err.Error())
			}
		}

		// Create session
		token, err := sessionInteractor.Set(created.ID, 24*time.Hour)
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

		writer.WriteResponseJSON(w, http.StatusCreated, created)
	}
}

// PutAvatar returns handler with environment which adds or updates profile avatar
// @Summary Put avatar
// @Description Update profile avatar
// @ID put-avatar
// @Accept multipart/form-data
// @Produce json
// @Param avatar body png false "Avatar"
// @Success 200 {object} profile.Avatar "Profile avatar is updated successfully"
// @Failure 400 "Incorrect request data"
// @Failure 403 "Not authorized"
// @Failure 500 "Database error"
// @Router /avatars [PUT]
func PutAvatar(profileInteractor *usecase.ProfileInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)
		id := r.Context().Value(middleware.ProfileID).(uint64)

		err := r.ParseMultipartForm(5 * (1 << 20)) // max size 5 MB
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Save new profile avatar
		avatar, header, err := r.FormFile("avatar")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer avatar.Close()
		filename := header.Filename
		dir := "upload/img/"

		avatarPath, err := saveAvatar(profileInteractor, avatar, filename, dir, id)
		if err != nil {
			log.Error(err.Error(),
				zap.Uint64("profile_id", id))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := &profile.Avatar{
			Path: avatarPath,
		}
		writer.WriteResponseJSON(w, http.StatusOK, result)
	}
}

// GetScoreboard returns handler with environment which processes request for getting score
// @Summary Get score
// @Description Get score
// @ID get-score
// @Produce json
// @Param limit query int false "limit number"
// @Param start query int false "start index"
// @Success 200 {object} profile.Profiles "Scoreboard found successfully"
// @Failure 400 "Incorrect request data"
// @Failure 500 "Database error"
// @Router /scores [GET]
func GetScoreboard(profileInteractor *usecase.ProfileInteractor) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)

		limit, limitErr := strconv.ParseInt(r.FormValue("limit"), 10, 64)
		offset, offsetErr := strconv.ParseInt(r.FormValue("start"), 10, 64)

		if limitErr != nil && offsetErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if limitErr == nil && limit > 25 {
			limit = 25
		}

		profiles := profile.Profiles{
			List: []profile.Info{},
		}

		// Get scores
		var err error
		profiles.List, err = profileInteractor.GetProfiles(limit, offset)
		if err != nil {
			log.Error(err.Error(),
				zap.Int64("limit", limit),
				zap.Int64("offset", offset))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get total number of scores
		profiles.Total, err = profileInteractor.GetProfileCount()
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteResponseJSON(w, http.StatusOK, profiles)
	}
}
