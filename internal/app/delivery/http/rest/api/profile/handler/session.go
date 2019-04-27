package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/response"

	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/session"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/middleware"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/http/rest/writer"
	"github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/domain/profile"

	profileService "github.com/go-park-mail-ru/2019_1_The-Recovery-Team/internal/app/delivery/grpc/service/profile"

	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes"

	"github.com/mailru/easyjson"

	"github.com/asaskevich/govalidator"
	"github.com/jackc/pgx"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// GetSession returns handler with environment whick gets profile id of current session
// @Summary Get session
// @Description Get profile id of authorized client
// @ID get-session
// @Produce json
// @Success 200 {object} profile.ID "Profile ID found successfully"
// @Failure 403 "Not authorized"
// @Router /sessions [GET]
func GetSession() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id := r.Context().Value(middleware.ProfileID)
		result := &profile.ID{
			Id: id.(uint64),
		}
		writer.WriteResponseJSON(w, http.StatusOK, result)
	}
}

// PostSession returns handler with environment which creates session
// @Summary Post session
// @Description Creates client session
// @ID post-session
// @Accept json
// @Param profile_login body profile.Login true "Email, password"
// @Success 200 {object} profile.Profile "Session is created successfully"
// @Failure 400 "Incorrect request data"
// @Failure 422 {object} handler.Error "Invalid request data"
// @Failure 403 "Not authorized"
// @Router /sessions [POST]
func PostSession(profileManager *profileService.ProfileClient, sessionManager *session.SessionClient) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)

		login := &profile.Login{}
		if err := easyjson.UnmarshalFromReader(r.Body, login); err != nil {
			r.Body.Close()
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.Body.Close()

		if isValid, err := govalidator.ValidateStruct(login); !isValid && err != nil {
			resp := response.Error{
				Description: err.Error(),
			}
			writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, resp)
			return
		}

		// Get profile by email and password
		request := &profileService.GetByEmailAndPasswordRequest{
			Email:    login.Email,
			Password: login.Password,
		}
		prof, err := (*profileManager).GetByEmailAndPassword(context.Background(), request)
		if err != nil {
			message := status.Convert(err).Message()
			if message == pgx.ErrNoRows.Error() {
				resp := response.Error{
					Description: "Invalid email or password",
				}
				writer.WriteResponseJSON(w, http.StatusUnprocessableEntity, resp)
				return
			}
			log.Error(message,
				zap.String("email", login.Email))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create session
		create := &session.Create{
			ProfileId: &session.ProfileId{
				Id: prof.Info.Id,
			},
			Expires: ptypes.DurationProto(24 * time.Hour),
		}
		sessionId, err := (*sessionManager).Set(context.Background(), create)
		if err != nil {
			log.Error(status.Convert(err).Message(),
				zap.Uint64("profile_id", prof.Info.Id))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionId.Id,
			Path:     ".",
			Expires:  time.Now().Add(24*time.Hour - 10*time.Minute),
			HttpOnly: true,
		})

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

// DeleteSession returns handler with environment which deletes session
// @Summary Delete session
// @Description Deletes client session
// @ID delete-session
// @Success 200 "Session is deleted successfully"
// @Failure 404 "Session not found"
// @Failure 403 "Not authorized"
// @Router /sessions [DELETE]
func DeleteSession(sessionManager *session.SessionClient) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		log := r.Context().Value("logger").(*zap.Logger)

		sessionID := &session.SessionId{
			Id: r.Context().Value(middleware.SessionID).(string),
		}
		if _, err := (*sessionManager).Delete(context.Background(), sessionID); err != nil {
			log.Error(status.Convert(err).Message(),
				zap.String("session_id", sessionID.Id))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
	}
}
