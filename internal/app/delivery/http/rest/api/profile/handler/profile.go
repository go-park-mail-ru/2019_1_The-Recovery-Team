package handler

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/middleware"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/response"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/saver"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/writer"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/profile"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/infrastructure/repository/postgresql"

	profileService "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"

	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes"

	"github.com/mailru/easyjson"

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

func saveAvatar(profileManager *profileService.ProfileClient, avatar multipart.File, filename, dir string, id uint64) (string, error) {
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
	request := &profileService.UpdateAvatarRequest{
		Id:     id,
		Avatar: avatarPath,
	}
	_, err = (*profileManager).UpdateAvatar(context.Background(), request)
	if err != nil {
		return "", errors.New(status.Convert(err).Message())
	}
	return avatarPath, nil
}

// List returns handler with environment which processes request for checking email or nickname existence
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
func GetProfiles(profileManager *profileService.ProfileClient) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)

		email := r.FormValue("email")
		nickname := r.FormValue("nickname")

		// Check email existence
		if email != "" && govalidator.IsEmail(email) {
			request := &profileService.GetByEmailRequest{
				Email: email,
			}
			_, err := (*profileManager).GetByEmail(context.Background(), request)
			if err != nil {
				message := status.Convert(err).Message()
				if message == pgx.ErrNoRows.Error() {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				log.Error(message,
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
			request := &profileService.GetByNicknameRequest{
				Nickname: nickname,
			}
			_, err := (*profileManager).GetByNickname(context.Background(), request)
			if err != nil {
				message := status.Convert(err).Message()
				if message == pgx.ErrNoRows.Error() {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				log.Error(message,
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

// Get returns handler with environment which processes request for getting profile by id
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
func GetProfile(profileManager *profileService.ProfileClient) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)

		id, err := strconv.ParseUint(ps.ByName("id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		request := &profileService.GetRequest{
			Id: id,
		}

		prof, err := (*profileManager).Get(context.Background(), request)
		if err != nil {
			message := status.Convert(err).Message()
			if message == pgx.ErrNoRows.Error() {
				log.Info(message,
					zap.Uint64("profile_id", id))
				w.WriteHeader(http.StatusNotFound)
				return
			}
			log.Error(message,
				zap.Uint64("profile_id", id))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		profileID := r.Context().Value(middleware.ProfileID)
		if profileID != id {
			prof.Email = ""
		}

		writer.WriteResponseJSON(w, http.StatusOK, &profile.Profile{
			Info: profile.Info{
				ID:       prof.Info.Id,
				Nickname: prof.Info.Nickname,
				Avatar:   prof.Info.Avatar,
				Score: profile.Score{
					Record: prof.Info.Score.Record,
					Win:    prof.Info.Score.Win,
					Loss:   prof.Info.Score.Loss,
				},
			},
			Email: prof.Email,
		})
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
func PutProfile(profileManager *profileService.ProfileClient) httprouter.Handle {
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
		if err := easyjson.UnmarshalFromReader(r.Body, data); err != nil {
			r.Body.Close()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.Body.Close()

		if isValid, err := govalidator.ValidateStruct(data); !isValid && err != nil {
			resp := response.Error{
				Description: err.Error(),
			}
			writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, resp)
			return
		}

		request := &profileService.UpdateRequest{
			Id:       id,
			Email:    data.Email,
			Nickname: data.Nickname,
		}
		if _, err = (*profileManager).Update(context.Background(), request); err != nil {
			message := status.Convert(err).Message()
			if message == EmailAlreadyExists || message == NicknameAlreadyExists {
				w.WriteHeader(http.StatusConflict)
				return
			}
			log.Error(message,
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
func PutProfilePassword(profileManager *profileService.ProfileClient) httprouter.Handle {
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
		if err := easyjson.UnmarshalFromReader(r.Body, data); err != nil {
			r.Body.Close()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.Body.Close()

		if isValid, err := govalidator.ValidateStruct(data); !isValid && err != nil {
			resp := response.Error{
				Description: err.Error(),
			}
			writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, resp)
			return
		}

		data.Password, err = postgresql.HashAndSalt(data.Password)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		request := &profileService.UpdatePasswordRequest{
			Id:          id,
			Password:    data.Password,
			PasswordOld: data.PasswordOld,
		}
		if _, err = (*profileManager).UpdatePassword(context.Background(), request); err != nil {
			message := status.Convert(err).Message()
			if message == IncorrectProfilePassword {
				resp := response.Error{
					Description: "Incorrect password",
				}
				writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, resp)
				return
			}
			log.Error(message,
				zap.Uint64("profile_id", id))
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
// @Failure 500 "Database error"`
// @Router /profiles [POST]
func PostProfile(profileManager *profileService.ProfileClient, sessionManager *session.SessionClient) httprouter.Handle {
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
			resp := response.Error{
				Description: err.Error(),
			}
			writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, resp)
			return
		}

		data.Password, err = postgresql.HashAndSalt(data.Password)
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create profile
		request := &profileService.CreateRequest{
			Email:    data.Email,
			Nickname: data.Nickname,
			Password: data.Password,
		}
		created, err := (*profileManager).Create(context.Background(), request)
		if err != nil {
			message := status.Convert(err).Message()
			if message == EmailAlreadyExists || message == NicknameAlreadyExists {
				log.Error(message,
					zap.String("email", data.Email),
					zap.String("nickname", data.Nickname))
				w.WriteHeader(http.StatusConflict)
				return
			}
			resp := response.Error{
				Description: message,
			}
			writer.WriteResponseJSON(w, http.StatusInternalServerError, resp)
			return
		}

		// Save profile avatar
		avatar, header, err := r.FormFile("avatar")
		if err == nil {
			defer avatar.Close()
			filename := header.Filename
			dir := "upload/img/"

			if avatarPath, err := saveAvatar(profileManager, avatar, filename, dir, created.Id); err == nil {
				created.Avatar = avatarPath
			} else {
				log.Warn(status.Convert(err).Message())
			}
		}

		// Create session
		create := &session.Create{
			ProfileId: &session.ProfileId{
				Id: created.Id,
			},
			Expires: ptypes.DurationProto(24 * time.Hour),
		}

		sessionId, err := (*sessionManager).Set(context.Background(), create)
		if err != nil {
			log.Error(status.Convert(err).Message(),
				zap.Uint64("created_id", created.Id))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionId.Id,
			Path:     "/",
			Expires:  time.Now().Add(24*time.Hour - 10*time.Minute),
			HttpOnly: true,
		})

		writer.WriteResponseJSON(w, http.StatusCreated, profile.Created{
			ID:       created.Id,
			Email:    created.Email,
			Nickname: created.Nickname,
			Avatar:   created.Avatar,
		})
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
func PutAvatar(profileManager *profileService.ProfileClient) httprouter.Handle {
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

		avatarPath, err := saveAvatar(profileManager, avatar, filename, dir, id)
		if err != nil {
			log.Error(status.Convert(err).Message(),
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
func GetScoreboard(profileManager *profileService.ProfileClient) httprouter.Handle {
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

		// Get scores
		request := &profileService.ListRequest{
			Limit:  limit,
			Offset: offset,
		}
		resp, err := (*profileManager).List(context.Background(), request)
		if err != nil {
			log.Error(status.Convert(err).Message(),
				zap.Int64("limit", limit),
				zap.Int64("offset", offset))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		profiles := profile.Profiles{
			List: []profile.Info{},
		}
		for _, info := range resp.List {
			profiles.List = append(profiles.List, profile.Info{
				ID:       info.Id,
				Nickname: info.Nickname,
				Avatar:   info.Avatar,
				Score: profile.Score{
					Record: info.Score.Record,
					Win:    info.Score.Win,
					Loss:   info.Score.Loss,
				},
			})
		}

		// Get total number of scores
		total, err := (*profileManager).Count(context.Background(), &profileService.Nothing{})
		if err != nil {
			log.Error(status.Convert(err).Message())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		profiles.Total = total.Count

		writer.WriteResponseJSON(w, http.StatusOK, profiles)
	}
}
