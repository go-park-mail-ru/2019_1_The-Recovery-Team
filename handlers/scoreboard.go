package handlers

import (
	"api/models"
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
func GetScoreboard(env *models.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, limitErr := strconv.ParseInt(r.FormValue("limit"), 10, 64)
		offset, offsetErr := strconv.ParseInt(r.FormValue("start"), 10, 64)

		if limitErr == nil && limit > 25 {
			limit = 25
		}

		result := models.Profiles{
			List: []models.Profile{},
		}
		var err error
		switch {
		case limitErr == nil && offsetErr == nil:
			{
				env.Dbm.FindAll(&result.List, QueryProfilesWithLimitAndOffset, limit, offset)
			}
		case limitErr == nil:
			{
				env.Dbm.FindAll(&result.List, QueryProfilesWithLimit, limit)
			}
		case offsetErr == nil:
			{
				env.Dbm.FindAll(&result.List, QueryProfilesWithOffset, offset)
			}
		default:
			{
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		env.Dbm.Find(&result.Total, QueryCountProfilesNumber)

		writeResponseJSON(w, http.StatusOK, result)
	}
}
