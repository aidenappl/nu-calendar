package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/aidenappl/nu-calendar/db"
	"github.com/aidenappl/nu-calendar/structs"
)

type QueryableFields struct {
	Email  *string
	UserID *int
	Name   *string
}

type GetUserRequest struct {
	UserID int
	Email  string
}

func GetUser(db db.Queryable, req GetUserRequest) (*structs.User, error) {
	var qFields = &QueryableFields{}

	if req.UserID != 0 {
		qFields.UserID = &req.UserID
	}

	if req.Email != "" {
		qFields.Email = &req.Email
	}

	users, err := ListUsers(db, ListUsersRequest{
		Limit:  1,
		Offset: 0,
		Query:  qFields,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	if len(*users) == 0 {
		return nil, nil
	}

	return &(*users)[0], nil
}

type ListUsersRequest struct {
	Limit  uint64
	Offset uint64
	Query  *QueryableFields
}

func ListUsers(db db.Queryable, req ListUsersRequest) (*[]structs.User, error) {
	q := sq.Select(
		"core.users.id",
		"core.users.name",
		"core.users.email",
		"core.users.inserted_at",
	).From("core.users")

	if req.Query != nil {
		if req.Query.Email != nil {
			q = q.Where(sq.Eq{"core.users.email": *req.Query.Email})
		}

		if req.Query.UserID != nil {
			q = q.Where(sq.Eq{"core.users.id": *req.Query.UserID})
		}

		if req.Query.Name != nil {
			q = q.Where(sq.Eq{"core.users.name": *req.Query.Name})
		}
	}

	q = q.Limit(req.Limit).Offset(req.Offset)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	defer rows.Close()

	var users []structs.User
	for rows.Next() {
		var user structs.User
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		users = append(users, user)
	}

	return &users, nil
}
