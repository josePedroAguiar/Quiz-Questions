package main
 
import (
    "fmt"
    "database/sql"
    _ "github.com/lib/pq"

)
 
const (
    host     = "localhost"
    port     = 5432
    userdb     = "postgres"
    password = "postgres"
    dbname   = "QUIZ_DB_1"
)
 
func connect() *sql.DB {
        // connection string
    psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, userdb, password, dbname)
         
        // open database
    db, err := sql.Open("postgres", psqlconn)
    CheckError(err)
     
    // close database
    
 
    // check db
    err = db.Ping()
    CheckError(err)
 
    fmt.Println("Connected!")
    return db;
}


func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}