package dto

var (
	REGISTRY_NOT_PROCESSED = "NOT_PROCESSED"
	INSERT_SENTENCE        = `INSERT INTO user_msg (message, status) VALUES (@Message, @Status );`
	CONNECTION_SUCCESFUL   = "Connected to SQL Server Database - ABC"
	USER_NULL              = "MigrateUser: db is null"
)
