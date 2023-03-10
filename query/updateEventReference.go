package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/aidenappl/nu-calendar/db"
)

type UpdateEventReferenceRequest struct {
	EventReferenceID int
	CustomSummary    *string
	LocationAddress  *string
	LocationBuilding *string
	LocationRoom     *string
	IncludeLocation  *bool
}

func UpdateEventReference(db db.Queryable, req UpdateEventReferenceRequest) error {
	q := sq.Update("nu_calendar.cal_reference_links").
		Where(sq.Eq{"nu_calendar.cal_reference_links.id": req.EventReferenceID})

	if req.LocationAddress != nil {
		q = q.Set("nu_calendar.cal_reference_links.location_address", *req.LocationAddress)
	}

	if req.LocationBuilding != nil {
		q = q.Set("nu_calendar.cal_reference_links.location_building", *req.LocationBuilding)
	}

	if req.LocationRoom != nil {
		q = q.Set("nu_calendar.cal_reference_links.location_room", *req.LocationRoom)
	}

	if req.CustomSummary != nil {
		q = q.Set("nu_calendar.cal_reference_links.custom_summary", *req.CustomSummary)
	}

	if req.IncludeLocation != nil {
		q = q.Set("nu_calendar.cal_reference_links.include_location", *req.IncludeLocation)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf("error building sql query: %v", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing sql query: %v", err)
	}

	return nil
}
