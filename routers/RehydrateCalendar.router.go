package routers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aidenappl/nu-calendar/db"
	"github.com/aidenappl/nu-calendar/errors"
	"github.com/aidenappl/nu-calendar/query"
)

func RehydrateCalendar(w http.ResponseWriter, r *http.Request) {
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

	// Get the courses for that EAB link
	events, err := GetUserCalendar(calendar.EabURI)
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

	// Remove duplicates within the slice
	classes = removeDuplicateStr(classes)

	// Remove any classes that already exist in the calendar
	for _, e := range calendar.EventReferences {
		for i, c := range classes {
			if c == e.Event {
				classes = append(classes[:i], classes[i+1:]...)
			}
		}
	}

	// Create the courses for that user
	for _, c := range classes {
		err = query.CreateEventReference(db.DB, query.CreateEventReferenceRequest{
			CalendarID: calendar.ID,
			Event:      c,
		})
		if err != nil {
			errors.SendError(w, "error creating event reference: "+err.Error(), "", http.StatusInternalServerError)
			return
		}
	}

	// Get the calendar
	calendar, err = query.GetCalendar(db.DB, query.GetCalendarRequest{
		CalendarID: &calendar.ID,
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
