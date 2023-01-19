package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/aidenappl/nu-calendar/db"
)

type GetEventReferencesRequest struct {
	CalendarID       *int
	EventReferenceID *int
}

func GetEventReferences(db db.Queryable, req GetEventReferencesRequest) (*[]EventReference, error) {
	q := sq.Select(
		"nu_calendar.cal_reference_links.id",
		"nu_calendar.cal_reference_links.cal_id",
		"nu_calendar.cal_reference_links.event",
		"nu_calendar.cal_reference_links.custom_summary",
		"nu_calendar.cal_reference_links.location_building",
		"nu_calendar.cal_reference_links.location_room",
		"nu_calendar.cal_reference_links.location_address",
		"nu_calendar.cal_reference_links.include_location",
		"nu_calendar.cal_reference_links.inserted_at",
	).
		From("nu_calendar.cal_reference_links")

	if req.CalendarID != nil {
		q = q.Where(sq.Eq{"nu_calendar.cal_reference_links.cal_id": *req.CalendarID})
	}

	if req.EventReferenceID != nil {
		q = q.Where(sq.Eq{"nu_calendar.cal_reference_links.id": *req.EventReferenceID})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building sql query: %v", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing sql query: %v", err)
	}

	var eventReferences []EventReference

	for rows.Next() {
		var eventReference EventReference

		err = rows.Scan(
			&eventReference.ID,
			&eventReference.CalID,
			&eventReference.Event,
			&eventReference.CustomSummary,
			&eventReference.LocationBuilding,
			&eventReference.LocationRoom,
			&eventReference.LocationAddress,
			&eventReference.IncludeLocation,
			&eventReference.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		eventReferences = append(eventReferences, eventReference)
	}

	return &eventReferences, nil
}

type GetEventReferenceRequest struct {
	EventReferenceID int
}

func GetEventReference(db db.Queryable, req GetEventReferenceRequest) (*EventReference, error) {
	eventReferences, err := GetEventReferences(db, GetEventReferencesRequest{EventReferenceID: &req.EventReferenceID})
	if err != nil {
		return nil, err
	}

	if len(*eventReferences) == 0 {
		return nil, nil
	}

	return &(*eventReferences)[0], nil
}
