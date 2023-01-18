package routers

import (
	"encoding/json"
	"net/http"

	"github.com/aidenappl/nu-calendar/db"
	"github.com/aidenappl/nu-calendar/errors"
	"github.com/aidenappl/nu-calendar/query"
)

func HandleGetCalendar(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get("slug")
	if slug == "" {
		errors.SendError(w, "missing required parameter: slug", "", http.StatusBadRequest)
		return
	}

	calendar, err := query.GetCalendar(db.DB, query.GetCalendarRequest{
		CalendarSlug: &slug,
	})
	if err != nil {
		errors.SendError(w, "error getting calendar", err.Error(), http.StatusInternalServerError)
		return
	}

	if calendar == nil {
		errors.SendError(w, "calendar not found", "", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(calendar)

}
