package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"users-backup-aws-lambda/pkg/dto"
)

var db *sql.DB

type (
	databaseConnection struct{}
)

func NewDatabaseConnection() *databaseConnection {
	return &databaseConnection{}
}

func (d *databaseConnection) getParamsConnection(connectionString string) *url.URL {

	var secret dto.Secret
	errorParsingSecret := json.Unmarshal([]byte(connectionString), &secret)

	if errorParsingSecret != nil {
		fmt.Println(errorParsingSecret)
	}

	portInteger, _ := strconv.Atoi(secret.Port)

	query := url.Values{}

	return &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(secret.User, secret.Password),
		Host:     fmt.Sprintf("%s:%d", secret.Ip, portInteger),
		RawQuery: query.Encode(),
	}

}

func (d *databaseConnection) Open(connectionString string) {

	var errorConnection error
	paramsConnection := d.getParamsConnection(connectionString)
	db, errorConnection = sql.Open(dto.SQL_ENGINE, paramsConnection.String())

	if errorConnection != nil {
		fmt.Println(errorConnection)
	}

	var err error

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf(dto.CONNECTION_SUCCESFUL)
}

func (*databaseConnection) MigrateUser(userToInsert string) (int64, error) {
	ctx := context.Background()
	var err error

	if db == nil {
		err = errors.New(dto.USER_NULL)
		return -1, err
	}

	// Check if storage is alive.
	err = db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	sqlSentence := dto.INSERT_SENTENCE

	stmt, err := db.Prepare(sqlSentence)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		sql.Named("Message", userToInsert),
		sql.Named("Status", dto.REGISTRY_NOT_PROCESSED))
	if err != nil {
		return -1, err
	}

	return 1, nil
}
