package nulltypes

import "database/sql"

func NewNullString(str *string) sql.NullString {
	n := sql.NullString{Valid: str != nil}
	if str != nil {
		n.String = *str
	}
	return n
}
