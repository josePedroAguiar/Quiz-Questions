package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "strings"

)

func getAll(db *sql.DB)[]Question {
    rows, err := db.Query(`SELECT * FROM tblQuestions`)
    CheckError(err)
    var questions []Question
    defer rows.Close()
    for rows.Next() { 
        var question Question
        err = rows.Scan( 
            &question.ID,
			&question.Question,
			&question.Description,
			&question.Answers.AnswerA,
			&question.Answers.AnswerB,
			&question.Answers.AnswerC,
			&question.Answers.AnswerD,
			&question.Answers.AnswerE,
			&question.Answers.AnswerF,
			&question.MultipleCorrectAnswers,
			&question.CorrectAnswers.AnswerACorrect,
			&question.CorrectAnswers.AnswerBCorrect,
			&question.CorrectAnswers.AnswerCCorrect,
			&question.CorrectAnswers.AnswerDCorrect,
			&question.CorrectAnswers.AnswerECorrect,
			&question.CorrectAnswers.AnswerFCorrect,
			&question.CorrectAnswer,
			&question.Explanation,
			&question.Tip,
			&question.Category,
			&question.Difficulty,)

      
        CheckError(err)
        question.Tags = getTagsFromString(question.Category, ";")
        questions = append(questions, question)
        
    }
    
    CheckError(err)
    return questions;


}



func getByDifficulty(db *sql.DB, diff string)[]Question {
    rows, err := db.Query(`SELECT * FROM tblQuestions WHERE difficulty = $1`,diff)
    CheckError(err)
    var questions []Question
    defer rows.Close()
    for rows.Next() { 
        var question Question
        err = rows.Scan( 
            &question.ID,
			&question.Question,
			&question.Description,
			&question.Answers.AnswerA,
			&question.Answers.AnswerB,
			&question.Answers.AnswerC,
			&question.Answers.AnswerD,
			&question.Answers.AnswerE,
			&question.Answers.AnswerF,
			&question.MultipleCorrectAnswers,
			&question.CorrectAnswers.AnswerACorrect,
			&question.CorrectAnswers.AnswerBCorrect,
			&question.CorrectAnswers.AnswerCCorrect,
			&question.CorrectAnswers.AnswerDCorrect,
			&question.CorrectAnswers.AnswerECorrect,
			&question.CorrectAnswers.AnswerFCorrect,
			&question.CorrectAnswer,
			&question.Explanation,
			&question.Tip,
			&question.Category,
			&question.Difficulty,)

      
        CheckError(err)
        question.Tags = getTagsFromString(question.Category, ";")
        questions = append(questions, question)
        
    }
    
    CheckError(err)
    return questions;
}

func getByCategory(db *sql.DB, category string)[]Question {
	str:= "%"+category+"%"
    rows, err := db.Query(`SELECT * FROM tblQuestions WHERE category LIKE $1`,str)
    CheckError(err)
    var questions []Question
    defer rows.Close()
    for rows.Next() { 
        var question Question
        err = rows.Scan( 
            &question.ID,
			&question.Question,
			&question.Description,
			&question.Answers.AnswerA,
			&question.Answers.AnswerB,
			&question.Answers.AnswerC,
			&question.Answers.AnswerD,
			&question.Answers.AnswerE,
			&question.Answers.AnswerF,
			&question.MultipleCorrectAnswers,
			&question.CorrectAnswers.AnswerACorrect,
			&question.CorrectAnswers.AnswerBCorrect,
			&question.CorrectAnswers.AnswerCCorrect,
			&question.CorrectAnswers.AnswerDCorrect,
			&question.CorrectAnswers.AnswerECorrect,
			&question.CorrectAnswers.AnswerFCorrect,
			&question.CorrectAnswer,
			&question.Explanation,
			&question.Tip,
			&question.Category,
			&question.Difficulty,)

      
        CheckError(err)
		question.Tags = getTagsFromString(question.Category, ";")
        questions = append(questions, question)
        
    }
    CheckError(err)
    return questions;
}
func getByCategoryAndDifficulty(db *sql.DB, category string,diff string)[]Question {
	str:= "%"+category+"%"
    rows, err := db.Query(`SELECT * FROM tblQuestions WHERE category LIKE $1 and difficulty = $2`,str,diff)
    CheckError(err)
    var questions []Question
    defer rows.Close()
    for rows.Next() { 
        var question Question
        err = rows.Scan( 
            &question.ID,
			&question.Question,
			&question.Description,
			&question.Answers.AnswerA,
			&question.Answers.AnswerB,
			&question.Answers.AnswerC,
			&question.Answers.AnswerD,
			&question.Answers.AnswerE,
			&question.Answers.AnswerF,
			&question.MultipleCorrectAnswers,
			&question.CorrectAnswers.AnswerACorrect,
			&question.CorrectAnswers.AnswerBCorrect,
			&question.CorrectAnswers.AnswerCCorrect,
			&question.CorrectAnswers.AnswerDCorrect,
			&question.CorrectAnswers.AnswerECorrect,
			&question.CorrectAnswers.AnswerFCorrect,
			&question.CorrectAnswer,
			&question.Explanation,
			&question.Tip,
			&question.Category,
			&question.Difficulty,)

      
        CheckError(err)
        question.Tags = getTagsFromString(question.Category, ";")
        questions = append(questions, question)
        
    }
    CheckError(err)
    return questions;
}



///
///
/// Auxiliary Functions
///
///

func getTagsFromString(tagNamesStr string, split string) []Tag {
	tagNames := strings.Split(tagNamesStr, split)
	tags := make([]Tag, 0, len(tagNames)) // Use a slice with initial capacity

	for i, tagName := range tagNames {
        
		if i == len(tagNames)-1 {
			break 
		}
		tags = append(tags, Tag{Name: tagName})
	}

	return tags
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
 
