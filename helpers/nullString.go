package helpers

import (
	"database/sql"

	"github.com/google/uuid"
)


func NewNullString(s *string) sql.NullString {
	if s != nil {
		return sql.NullString{
			String: *s,
			Valid: true,
		}
	}

	return sql.NullString{
		String: "",
		Valid: false,
	}
}

func NewNullInt(i *int) sql.NullInt32 {
	if i != nil {
		return sql.NullInt32{
			Int32: int32(*i),
			Valid: true,
		}
	}

	return sql.NullInt32{
		Int32: 0,
		Valid: false,
	}
}

func NewNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u != nil {
		return uuid.NullUUID{
			UUID: *u,
			Valid:  true,
		}
	}

	return uuid.NullUUID{
		UUID: uuid.UUID{},
		Valid:  false,
	}
}