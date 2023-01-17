package routers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aidenappl/nu-calendar/errors"
	ics "github.com/arran4/golang-ical"
)

type HandleLaunchpadRequest struct {
	CalendarID int `json:"calendar_id"`
}

func HandleLaunchpad(w http.ResponseWriter, r *http.Request) {

	events, err := GetUserCalendar("https://northeastern.campus.eab.com/cal/OQbi8N5YcCWN/GradesFirst.ics")
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

	json.NewEncoder(w).Encode(classes)
}

func GetUserCalendar(CalendarURI string) (*[]Event, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, CalendarURI, nil)

	if err != nil {
		return nil, fmt.Errorf("error creating http request: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making http request: %v", err)
	}
	defer res.Body.Close()

	c, err := ics.ParseCalendar(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error parsing calendar: %v", err)
	}

	var events []Event

	for _, e := range c.Components {
		var event Event
		for _, v := range e.UnknownPropertiesIANAProperties() {
			if v.IANAToken == "SUMMARY" {
				event.Summary = v.Value
			}
			if v.IANAToken == "DTSTART" {
				event.StartTime = v.Value
			}
			if v.IANAToken == "DTEND" {
				event.EndTime = v.Value
			}
			if v.IANAToken == "UID" {
				event.UID = v.Value
			}
			if v.IANAToken == "URL" {
				event.URL = v.Value
			}

		}
		events = append(events, event)
	}

	return &events, nil
}
