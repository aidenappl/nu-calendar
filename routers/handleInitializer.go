package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aidenappl/nu-calendar/db"
	"github.com/aidenappl/nu-calendar/errors"
	"github.com/aidenappl/nu-calendar/query"
)

type HandleInitializerRequest struct {
	NortheasternEmail string `json:"northeastern_email"`
	EABUri            string `json:"eab_uri"`
}

func HandleInitializer(w http.ResponseWriter, r *http.Request) {
	var body HandleInitializerRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		errors.SendError(w, "error decoding request body: "+err.Error(), "", http.StatusBadRequest)
		return
	}

	if body.NortheasternEmail == "" {
		errors.SendError(w, "missing or empty field: northeastern_email", "", http.StatusBadRequest)
		return
	}

	if body.EABUri == "" {
		errors.SendError(w, "missing or empty field: eab_uri", "", http.StatusBadRequest)
		return
	}

	if !strings.Contains(body.EABUri, "https://northeastern.campus.eab.com") {
		errors.SendError(w, "invalid eab uri", "", http.StatusBadRequest)
		return
	}

	// Check if the user already has a calendar
	calendar, err := query.GetCalendar(db.DB, query.GetCalendarRequest{
		NortheasternEmail: &body.NortheasternEmail,
	})
	if err != nil {
		errors.SendError(w, "error getting calendar: "+err.Error(), "", http.StatusInternalServerError)
		return
	}

	if calendar != nil {
		errors.SendError(w, "calendar already exists", "", http.StatusBadRequest)
		return
	}

	// Create the calendar link
	calID, err := query.CreateCalLink(db.DB, query.CreateCalLinkRequest{
		NortheasternEmail: body.NortheasternEmail,
		EABUri:            body.EABUri,
	})
	if err != nil {
		errors.SendError(w, "error creating calendar link: "+err.Error(), "", http.StatusInternalServerError)
		return
	}

	// Get the courses for that EAB link
	events, err := GetUserCalendar(body.EABUri)
	if err != nil {
		fmt.Println(err)
		return
	}

	if events == nil {
		errors.SendError(w, "no events found", "", http.StatusBadRequest)
		return
	}

	classes := []string{}
	for _, e := range *events {
		classes = append(classes, e.Summary)
	}

	classes = removeDuplicateStr(classes)

	// Create the courses for that user
	for _, c := range classes {
		err = query.CreateEventReference(db.DB, query.CreateEventReferenceRequest{
			CalendarID: *calID,
			Event:      c,
		})
		if err != nil {
			errors.SendError(w, "error creating event reference: "+err.Error(), "", http.StatusInternalServerError)
			return
		}
	}

	// Get the calendar
	calendar, err = query.GetCalendar(db.DB, query.GetCalendarRequest{
		CalendarID: calID,
	})
	if err != nil {
		errors.SendError(w, "error getting calendar:"+err.Error(), "", http.StatusInternalServerError)
		return
	}

	if calendar == nil {
		errors.SendError(w, "calendar not found", "", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(calendar)
}
