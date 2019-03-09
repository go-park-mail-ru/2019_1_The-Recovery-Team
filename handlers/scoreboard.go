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
// @Failure 500 "Database error"
// @Router /scores [GET]
func GetScoreboard(env *models.Env) http.HandlerFunc {
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
			return
		}

		writeResponseJSON(w, http.StatusOK, result)
	}
}
