
package main

import (
	"net/http"
    "encoding/json"
	"fmt"
	"log"
	"os/user"
    "io/ioutil"
    "net/url"
    "os"
	"time"
	//"github.com/gin-gonic/gin"
)

type Question struct {
	ID                      int    `json:"id"`
	Question                string `json:"question"`
	Description             string `json:"description"`
	Answers                 Answers `json:"answers"`
	MultipleCorrectAnswers  string  `json:"multiple_correct_answers"`
	CorrectAnswers          CorrectAnswers `json:"correct_answers"`
	CorrectAnswer           string `json:"correct_answer"`
	Explanation             string `json:"explanation"`
	Tip                     string `json:"tip"`
	Tags                    []Tag `json:"tags"`
	Category                string `json:"category"`
	Difficulty              string `json:"difficulty"`
}

type Answers struct {
	AnswerA string `json:"answer_a"`
	AnswerB string `json:"answer_b"`
	AnswerC string `json:"answer_c"`
	AnswerD string `json:"answer_d"`
	AnswerE string `json:"answer_e"`
	AnswerF string `json:"answer_f"`
}

type CorrectAnswers struct {
	AnswerACorrect string `json:"answer_a_correct"`
	AnswerBCorrect string `json:"answer_b_correct"`
	AnswerCCorrect string `json:"answer_c_correct"`
	AnswerDCorrect string `json:"answer_d_correct"`
	AnswerECorrect string `json:"answer_e_correct"`
	AnswerFCorrect string `json:"answer_f_correct"`
}

type Tag struct {
	Name string `json:"name"`
}






// Hello returns a greeting for the named person.
func Hello(name string) string {
    // Return a greeting that embeds the name in a message.
    message := fmt.Sprintf("Hi, %v. Welcome!", name)
    return message
}


func api(category string,difficulty string) []Question {
	params := url.Values{
		"apiKey":     {"7XQrCvRqgHNM9XvEultYNAdlJVIY0QJMf4htCts3"},
		"limit":      {"30"},
		"tag":   {category},
		"difficulty": {difficulty},
	}

	url := "https://quizapi.io/api/v1/questions?" + params.Encode()

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var questions []Question

	json.Unmarshal([]byte(responseData), &questions)


	/*for _, question := range questions {
		fmt.Println("ID:", question.ID)
		fmt.Println("Question:", question.Question)
		fmt.Println("Answers:")
		fmt.Println("  A:", question.Answers.AnswerA)
		fmt.Println("  B:", question.Answers.AnswerB)
		fmt.Println("  C:", question.Answers.AnswerC)
		fmt.Println("  D:", question.Answers.AnswerD)
		fmt.Println("  E:", question.Answers.AnswerE)
		fmt.Println("  F:", question.Answers.AnswerF)
		fmt.Println("Correct Answer:", question.CorrectAnswer)
		fmt.Println("Explanation:", question.Explanation)
		fmt.Println("Tags:")
		for _, tag := range question.Tags {
			fmt.Println("  -", tag.Name)
		}
		fmt.Println("Category:", question.Category)
		fmt.Println("Difficulty:", question.Difficulty)
		fmt.Println()
	}*/

   return questions
	//fmt.Println(string(responseData))
}





func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	username := user.Name


	diffs := [3]string{"Easy","Medium","Hard"}
	category := []string{"Linux","DevOps","Networking","PHP", "JavaScript", "Python","Cloud","Docker","Kubernetes"}
	
	/*router := gin.Default()
    router.Run("localhost:8081")*/
	for 1==1 { //probably is more  elegant and clever approach
		for _,diff := range diffs{
			for _,cat := range category{
				questions:=api(cat,diff)
				db:=connect()
				insert(db,questions)
				time.Sleep(1 * time.Second)
		}
			time.Sleep(2 * time.Second)
		}
	fmt.Println(Hello(username))
	//time.Sleep(1 * time.Hour)
	time.Sleep(30 * time.Second)
	}
}
