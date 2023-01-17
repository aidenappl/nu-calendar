package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/aidenappl/nu-calendar/db"
	"github.com/aidenappl/nu-calendar/tools"
)

type CreateCalLinkRequest struct {
	NortheasternEmail string
	EABUri            string
}

func CreateCalLink(db db.Queryable, req CreateCalLinkRequest) (*int, error) {

	generatedSlug := tools.RandomString(8, tools.LetterNumberBytes)

	query, args, err := sq.Insert("nu_calendar.cal_links").
		Columns("nu_email", "eab_uri", "slug").
		Values(req.NortheasternEmail, req.EABUri, generatedSlug).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building sql query: %v", err)
	}

	rows, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing sql query: %v", err)
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %v", err)
	}

	var intID = int(id)
	return &intID, nil
}
