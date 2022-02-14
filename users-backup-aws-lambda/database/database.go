package database

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
)

var db *sql.DB

const NOT_PROCESSED string = "NOT PROCESSED"

type Secret struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Db       string `json:"db"`
	Server   string `json:"server"`
}

type (
	databaseConnection struct{}
)

func NewDatabaseConnection() *databaseConnection {
	return &databaseConnection{}
}

func (*databaseConnection) Open(connectionString string) {

	var secret Secret

	errorParsingSecret := json.Unmarshal([]byte(connectionString), &secret)

	if errorParsingSecret != nil {
		fmt.Println(errorParsingSecret)
	} else {
		fmt.Println(secret.Db, secret.Ip)
	}

	portInteger, _ := strconv.Atoi(secret.Port)

	query := url.Values{}
	query.Add("app name", "MyAppName")

	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(secret.User, secret.Password),
		Host:     fmt.Sprintf("%s:%d", secret.Ip, portInteger),
		RawQuery: query.Encode(),
	}
	db, errorConnection := sql.Open("sqlserver", u.String())

	if errorConnection != nil {
		fmt.Println(errorConnection)
	}

	var err error

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
}

func (*databaseConnection) MigrateUser(userToInsert string) (int64, error) {
	ctx := context.Background()
	var err error

	if db == nil {
		err = errors.New("MigrateUser: db is null")
		return -1, err
	}

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	sqlSentence := `INSERT INTO user_msg (message, status) VALUES (@Message, @Status );`

	stmt, err := db.Prepare(sqlSentence)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Message", userToInsert),
		sql.Named("Status", NOT_PROCESSED))
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}

	return newID, nil
}
