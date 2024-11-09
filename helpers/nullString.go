package helpers

import "database/sql"


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