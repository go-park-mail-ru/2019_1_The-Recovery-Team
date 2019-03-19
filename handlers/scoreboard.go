package handlers

import (
	"api/environment"
	"api/models"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// GetScoreboard returns handler with environment which processes request for getting scoreboard
// @Summary Get scoreboard
// @Description Get scoreboard
// @ID get-scoreboard
// @Produce json
// @Param limit query int false "limit number"
// @Param start query int false "start index"
// @Success 200 {object} models.Profiles "Scoreboard found successfully"
// @Failure 400 "Incorrect request data"
// @Failure 500 "Database error"
// @Router /scores [GET]
func GetScoreboard(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		profiles := models.Profiles{
			List: []models.ProfileInfo{},
		}

		var err error
		profiles.List, err = env.Dbm.GetProfiles(limit, offset)
		if err != nil {
			log.Error(err.Error(),
				zap.Int64("limit", limit),
				zap.Int64("offset", offset))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		profiles.Total, err = env.Dbm.GetProfilesNumber()
		if err != nil {
			log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		writeResponseJSON(w, http.StatusOK, profiles)
	}
}
