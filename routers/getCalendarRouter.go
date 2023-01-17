package routers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aidenappl/nu-calendar/db"
	"github.com/aidenappl/nu-calendar/errors"
	"github.com/aidenappl/nu-calendar/query"
	ics "github.com/arran4/golang-ical"
)

type Event struct {
	Summary   string
	StartTime string
	EndTime   string
	UID       string
	URL       string
}

type EMap struct {
	LocationBuilding *string
	LocationRoom     *string
	LocationAddress  *string
}

func GetCalendar(w http.ResponseWriter, r *http.Request) {
	CalSlug := r.URL.Query().Get("cal_slug")
	if CalSlug == "" {
		errors.SendError(w, "missing or empty field: cal_slug", "", http.StatusBadRequest)
		return
	}

	// Get the calendar
	calendar, err := query.GetCalendar(db.DB, query.GetCalendarRequest{
		CalendarSlug: &CalSlug,
	})
	if err != nil {
		errors.SendError(w, "error getting calendar:"+err.Error(), "", http.StatusInternalServerError)
		return
	}

	if calendar == nil {
		errors.SendError(w, "calendar not found", "", http.StatusBadRequest)
		return
	}

	var CoursesAndLocations = map[string]EMap{}

	// Convert the calendar fields to a string map
	for _, eventRef := range calendar.EventReferences {
		var courseAndLocation = EMap{}
		courseAndLocation.LocationBuilding = eventRef.LocationBuilding
		courseAndLocation.LocationRoom = eventRef.LocationRoom
		courseAndLocation.LocationAddress = eventRef.LocationAddress
		CoursesAndLocations[eventRef.Event] = courseAndLocation
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, calendar.EabURI, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	c, err := ics.ParseCalendar(res.Body)
	if err != nil {
		fmt.Println(err)
		return
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

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	cal.SetCalscale("GREGORIAN")
	cal.SetName("NU Timetable")
	cal.SetDescription("NU Timetable")
	cal.SetXWRCalName("NU Timetable")
	cal.SetXWRCalDesc("NU Timetable")

	for _, e := range events {

		event := cal.AddEvent(e.UID)

		layout := "20060102T150405Z" // format of the input string
		startTime, err := time.Parse(layout, e.StartTime)
		if err != nil {
			fmt.Println(err)
		}

		endTime, err := time.Parse(layout, e.EndTime)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		event.SetStartAt(startTime.In(time.Local))
		event.SetDtStampTime(time.Now())
		event.SetEndAt(endTime.In(time.Local))

		event.SetOrganizer("Northeastern University")
		event.SetURL(e.URL)
		event.SetClass("PUBLIC")
		event.SetSequence(0)

		var referenceData = CoursesAndLocations[e.Summary]

		if referenceData.LocationAddress == nil {
			var temp = ""
			referenceData.LocationAddress = &temp
		}

		if referenceData.LocationBuilding == nil {
			var temp = ""
			referenceData.LocationBuilding = &temp
		}

		if referenceData.LocationRoom == nil {
			var temp = ""
			referenceData.LocationRoom = &temp
		}

		event.SetLocation(fmt.Sprintf("%s - %s", *referenceData.LocationBuilding, *referenceData.LocationRoom))

		event.SetSummary(fmt.Sprintf("%s (%s - %s)", e.Summary, *referenceData.LocationBuilding, *referenceData.LocationRoom))

	}

	serializedCal := cal.Serialize()

	dtx := []byte(serializedCal)

	w.Header().Set("Content-Type", "text/calendar")
	w.WriteHeader(http.StatusOK)
	w.Write(dtx)
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
