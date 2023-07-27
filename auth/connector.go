package main
 
import (
    "fmt"
    "database/sql"
    _"github.com/lib/pq"

)
 
const (
    host     = "localhost"
    port     = 5433
    userdb     = "postgres"
    password = "postgres"
    dbname   = "QUIZ_AUTH"
)
 
func connect() *sql.DB {
    // Connection string
    psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, userdb, password, dbname)

    // Open database
    db, err := sql.Open("postgres", psqlconn)
    CheckError("sql.Open", err) // Provide context information for sql.Open error

    // Check database connection
    err = db.Ping()
    CheckError("db.Ping", err) // Provide context information for db.Ping error

    fmt.Println("Connected!")
    return db
}



func CheckError(context string, err error) {
    if err != nil {
        panic(fmt.Errorf("%s: %w", context, err))
    }
}

