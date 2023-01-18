package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/aidenappl/nu-calendar/db"
)

type GetCalendarRequest struct {
	CalendarID        *int
	CalendarSlug      *string
	NortheasternEmail *string
}

type Calendar struct {
	ID              int              `json:"id"`
	Nickname        *string          `json:"nickname"`
	NuEmail         string           `json:"nu_email"`
	EabURI          string           `json:"eab_uri"`
	Slug            string           `json:"slug"`
	EventReferences []EventReference `json:"event_references"`
	InsertedAt      string           `json:"inserted_at"`
}

type EventReference struct {
	ID               int     `json:"id"`
	CalID            int     `json:"cal_id"`
	Event            string  `json:"event"`
	LocationBuilding *string `json:"location_building"`
	LocationRoom     *string `json:"location_room"`
	LocationAddress  *string `json:"location_address"`
	InsertedAt       string  `json:"inserted_at"`
}

func GetCalendar(db db.Queryable, req GetCalendarRequest) (*Calendar, error) {
	q := sq.Select(
		"nu_calendar.cal_links.id",
		"nu_calendar.cal_links.nickname",
		"nu_calendar.cal_links.nu_email",
		"nu_calendar.cal_links.eab_uri",
		"nu_calendar.cal_links.slug",
		"nu_calendar.cal_links.inserted_at",
	).
		From("nu_calendar.cal_links").
		Limit(1).
		Offset(0)

	if req.CalendarID != nil {
		q = q.Where(sq.Eq{"nu_calendar.cal_links.id": *req.CalendarID})
	}

	if req.CalendarSlug != nil {
		q = q.Where(sq.Eq{"nu_calendar.cal_links.slug": *req.CalendarSlug})
	}

	if req.NortheasternEmail != nil {
		q = q.Where(sq.Eq{"nu_calendar.cal_links.nu_email": *req.NortheasternEmail})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building sql query: %v", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing sql query: %v", err)
	}

	defer rows.Close()

	var calendar Calendar

	for rows.Next() {
		err = rows.Scan(
			&calendar.ID,
			&calendar.Nickname,
			&calendar.NuEmail,
			&calendar.EabURI,
			&calendar.Slug,
			&calendar.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		references, err := GetEventReferences(db, GetEventReferencesRequest{CalendarID: &calendar.ID})
		if err != nil {
			return nil, fmt.Errorf("error getting event references: %v", err)
		}

		calendar.EventReferences = *references
	}

	if calendar.ID == 0 {
		return nil, nil
	}

	return &calendar, nil
}
