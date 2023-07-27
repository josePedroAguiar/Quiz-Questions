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


func insert(db *sql.DB,questions[]Question) {
    // func insert
    insertDynStmt :=`insert into tblQuestions ("id","question","description","answer_a","answer_b","answer_c","answer_d","answer_e","answer_f","multiple_correct_answers","answer_a_correct","answer_b_correct","answer_c_correct","answer_d_correct","answer_e_correct","answer_f_correct","correct_answer","explanation","tip","category","difficulty") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,$21)`
    //var question Question;
    var err error;
    for _, question := range questions {
         _, err = db.Exec(insertDynStmt,
        question.ID,
        question.Question,
        question.Description,
        question.Answers.AnswerA,
        question.Answers.AnswerB,
        question.Answers.AnswerC,
        question.Answers.AnswerD,
        question.Answers.AnswerE,
        question.Answers.AnswerF,
        question.MultipleCorrectAnswers,
        convertBool(question.CorrectAnswers.AnswerACorrect),
        convertBool(question.CorrectAnswers.AnswerBCorrect),
        convertBool(question.CorrectAnswers.AnswerCCorrect),
        convertBool(question.CorrectAnswers.AnswerDCorrect),
        convertBool(question.CorrectAnswers.AnswerECorrect),
        convertBool(question.CorrectAnswers.AnswerFCorrect),
        question.CorrectAnswer,
        question.Explanation,
        question.Tip,
        convertToString(question.Tags),
        question.Difficulty,
    )
    CheckError(err)
    defer db.Close()
}
}

func convertBool(str string) bool{
    return str == "true"


}
func convertToString(tags [] Tag) string{
    var finalstr string;
    for _, tag := range tags{
        finalstr=finalstr+tag.Name+";";
    }
    return finalstr;

}
 
func CheckError(err error) {
    if err != nil {
        fmt.Println(err)
    }
}