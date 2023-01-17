package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/aidenappl/nu-calendar/db"
)

type CreateEventReferenceRequest struct {
	CalendarID       int
	Event            string
	LocationBuilding *string
	LocationRoom     *string
	LocationAddress  *string
}

func CreateEventReference(db db.Queryable, req CreateEventReferenceRequest) error {

	var cols = []string{
		"cal_id",
		"event",
	}
	var vals = []interface{}{
		req.CalendarID,
		req.Event,
	}

	if req.LocationBuilding != nil {
		cols = append(cols, "location_building")
		vals = append(vals, req.LocationBuilding)
	}

	if req.LocationRoom != nil {
		cols = append(cols, "location_room")
		vals = append(vals, req.LocationRoom)
	}

	if req.LocationAddress != nil {
		cols = append(cols, "location_address")
		vals = append(vals, req.LocationAddress)
	}

	query, args, err := sq.Insert("nu_calendar.cal_reference_links").
		Columns(cols...).
		Values(vals...).
		ToSql()

	if err != nil {
		return fmt.Errorf("error building sql query: %v", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing sql query: %v", err)
	}

	return nil
}
