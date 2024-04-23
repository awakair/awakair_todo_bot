package dbqueries

type DbQuery = string

type UserFilled struct {
	Full         DbQuery
	LanguageCode DbQuery
	UtcOffset    DbQuery
	Empty        DbQuery
}

var InsertUser = UserFilled{
	Full: `INSERT INTO users (id, language_code, utc_offset)
	VALUES ($1::bigint, $2, $3);`,
	LanguageCode: `INSERT INTO users (id, language_code)
	VALUES ($1::bigint, $2);`,
	UtcOffset: `INSERT INTO users (id, utc_offset)
	VALUES ($1::bigint, $2);`,
	Empty: `INSERT INTO users (id)
	VALUES ($1::bigint);`,
}

var UpdateUser = UserFilled{
	Full: `UPDATE users
	SET language_code = $2, utc_offset = $3
	WHERE id = $1`,
	LanguageCode: `UPDATE users
	SET language_code = $2
	WHERE id = $1`,
	UtcOffset: `UPDATE users
	SET utc_offset = $3
	WHERE id = $1`,
}
