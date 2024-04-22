package dbqueries

type DbQuery = string

const (
	InsertNewUser DbQuery = `INSERT INTO users (id, language_code, utc_offset)
	VALUES ($1::bigint, $2, $3);`

	InsertNewUserLanguageCode DbQuery = `INSERT INTO users (id, language_code)
	VALUES ($1::bigint, $2);`

	InsertNewUserUtcOffset DbQuery = `INSERT INTO users (id, utc_offset)
	VALUES ($1::bigint, $2);`

	InsertNewUserEmpty = `INSERT INTO users (id)
	VALUES ($1::bigint);`
)
