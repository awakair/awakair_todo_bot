package dbqueries

type DbQuery = string

const (
	InsertNewUser DbQuery = "INSERT INTO users (id, language_code, utc_offset) VALUES ($1::bigint, $2, $3);"
)
