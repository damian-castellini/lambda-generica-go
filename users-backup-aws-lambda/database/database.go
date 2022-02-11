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

	// connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
	//	secret.Ip, secret.User, secret.Password, secret.Port, secret.Db)

	port, _ := strconv.Atoi(secret.Port)

	query := url.Values{}
	query.Add("app name", "MyAppName")

	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(secret.User, secret.Password),
		Host:   fmt.Sprintf("%s:%d", secret.Ip, port),
		// Path:  instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	db, errorConnection := sql.Open("sqlserver", u.String())

	if errorConnection != nil {
		fmt.Println(errorConnection)
	}

	var err error

	// Create connection pool
	// db, err = sql.Open("sqlserver", connString)
	// if err != nil {
	//		log.Fatal("Error creating connection pool: ", err.Error())
	//	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!\n")
}

func CreateEmployee(name string, location string) (int64, error) {
	ctx := context.Background()
	var err error

	if db == nil {
		err = errors.New("CreateEmployee: db is null")
		return -1, err
	}

	// Check if database is alive.
	err = db.PingContext(ctx)
	if err != nil {
		return -1, err
	}

	tsql := `INSERT INTO TestSchema.Employees (Name, Location) VALUES (@Name, @Location);
      select isNull(SCOPE_IDENTITY(), -1);`

	stmt, err := db.Prepare(tsql)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("Name", name),
		sql.Named("Location", location))
	var newID int64
	err = row.Scan(&newID)
	if err != nil {
		return -1, err
	}

	return newID, nil
}
