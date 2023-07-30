
package main

import (
	"net/http"
	"fmt"
	//"math/rand"
	"strconv"
	"net/url"
    "os"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
    "encoding/json"
    "io/ioutil"

	"strings"
	"math/rand"
	"time"
	"database/sql"

)

type Quiz struct {
	ID         int64 
	CreatedTime time.Time
	Name       string 
	UserID     int `json:"id"`
	Questions  []Question `json:"questions"` // One-to-many relationship with questions
}

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

type QuizUserAssociation struct {
	QuizID  int   `json:"quiz_id"`
	UserIDs []int `json:"user_ids"`
}



var MySigningKey = []byte(os.Getenv ("SECRET_KEY"))// Replace this with your secret key
//var MySigningKey = "ola"

func PostgetQuiz(c *gin.Context) {
	//db:=connect()
	// Define a struct to hold the data sent in the POST request
	type FormData struct {
		Difficulty string `form:"difficulty"`
		Category        string `form:"category"`
		Number     *int    `form:"number,default=-1"`
	}


	var formData FormData

	// Bind the POST data to the struct
	if err := c.ShouldBind(&formData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the required fields are empty


	if formData.Number!=nil&& *formData.Number <= 0{
		c.JSON(http.StatusOK, gin.H{"error": "Number should be bigger then 0"})
		return
	}


}




func getQuiz(c *gin.Context) {
	// Access query parameters using c.Query()
	difficulty := c.Query("difficulty")
	categorys := c.Query("categorys")
	numberStr := c.DefaultQuery("number", "all") // Default to 0 if the "number" query parameter is not provided
	var number int
	var err error
	// Convert numberStr to an integer
	if numberStr!="all"{
		number, err = strconv.Atoi(numberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number parameter"})
		return
	}
	
		if number <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Number should be bigger than 0"})
			return
		}
	}else{
		number=3
	}
	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", "http://localhost:8080/questions?", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	

	// Set the Authorization header
	req.Header.Set("Authorization", c.GetHeader("Authorization"))


	cats:=getTagsFromString(categorys,",")
	
	var questions []Question
	for _,category := range cats {
		var qs []Question
		// Set the query parameters
		params := url.Values{}
		number++
		params.Add("number", "2")
		params.Add("category", category.Name)
		params.Add("difficulty", difficulty)


		// Add the query parameters to the request URL
		req.URL.RawQuery = params.Encode()

		// Send the request
		client := http.Client{}
		resp, err := client.Do(req)
		
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}
		json.Unmarshal([]byte(body), &qs)
		questions = append(questions, qs...)
	}

	questions =  getRandomElements(questions, number) 

	//db := connect()


	id, _ := c.Get("id")
	username, _ := c.Get("username")
	email, _ := c.Get("email")
	isAdmin, _ := c.Get("is_admin")

	// Respond with the user and the questions as JSON
	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"user":     username,
		"email":    email,
		"is_admin": isAdmin,
		"questions": questions,
	})

}


func insertQuiz(db *sql.DB, quiz Quiz) error {
	_, err := db.Exec("INSERT INTO tblquizzes (id,name, tblusers_id) VALUES ($1, $2, $3)",quiz.ID, quiz.Name, quiz.UserID)
	return err
}

// Function to insert a question into the tblquizzes_questions table
func insertQuestion(db *sql.DB,quiz Quiz,i int) error {
	_, err := db.Exec("INSERT INTO tblquizzes_tblquestions (tblquizzes_id, questions_id,question,description,answer_a,answer_b,answer_c,answer_d,answer_e,answer_f,multiple_correct_answers,answer_a_correct,answer_b_correct,answer_c_correct,answer_d_correct,answer_e_correct,answer_f_correct,correct_answer,explanation,tip,category,difficulty) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)",
		quiz.ID, quiz.Questions[i].ID,
		quiz.Questions[i].Question,
		quiz.Questions[i].Description,
		quiz.Questions[i].Answers.AnswerA,
		quiz.Questions[i].Answers.AnswerB,
		quiz.Questions[i].Answers.AnswerC,
		quiz.Questions[i].Answers.AnswerD,
		quiz.Questions[i].Answers.AnswerE,
		quiz.Questions[i].Answers.AnswerF,
		quiz.Questions[i].MultipleCorrectAnswers,
			convertBool(quiz.Questions[i].CorrectAnswers.AnswerACorrect),
			convertBool(quiz.Questions[i].CorrectAnswers.AnswerBCorrect),
			convertBool(quiz.Questions[i].CorrectAnswers.AnswerCCorrect),
			convertBool(quiz.Questions[i].CorrectAnswers.AnswerDCorrect),
			convertBool(quiz.Questions[i].CorrectAnswers.AnswerECorrect),
			convertBool(quiz.Questions[i].CorrectAnswers.AnswerFCorrect),
			quiz.Questions[i].CorrectAnswer,
			quiz.Questions[i].Explanation,
			quiz.Questions[i].Tip,
			convertToString(quiz.Questions[i].Tags),
			quiz.Questions[i].Difficulty,
		
	)
	return err
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


func saveQuiz(c *gin.Context) {
	var quiz Quiz
 	db:=connect()
	// Read the request body
	body, err := ioutil.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error reading request body"})
		return
	}

	// Unmarshal the request body into the 'quiz' variable
	err = json.Unmarshal(body, &quiz)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error unmarshalling request body"})
		return
	}
	quiz.ID= generateUniqueSeed()
	c.JSON(http.StatusBadRequest,insertQuiz(db,quiz))

	for i := 0; i < len(quiz.Questions); i++  {
		c.JSON(http.StatusOK,insertQuestion(db,quiz,i))}
}


func getQuizAssign(c *gin.Context) {
	db := connect()
	// Get the user ID from the context and convert it to an integer
	userIDRaw, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
		return
	}
	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var quizzes []Quiz

	// Query the tblquizzes_tblusers table to get the quiz IDs associated with the user's ID
	rows, err := db.Query("SELECT tblquizzes_id FROM tblquizzes_tblusers WHERE tblusers_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quizzes 1"})
		return
	}
	defer rows.Close()

	// Iterate through the rows and fetch the quiz details for each quiz ID
	for rows.Next() {
		var quizID int
		if err := rows.Scan(&quizID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quiz details"})
			return
		}

		// Query the tblquizzes table to fetch the quiz details
		var quiz Quiz
		err := db.QueryRow("SELECT id, created_time, name,tblquizzes FROM tblquizzes WHERE id = ?", quizID).Scan(&quiz.ID, &quiz.CreatedTime, &quiz.Name,&quiz.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quiz details"})
			return
		}

		// Append the quiz to the quizzes slice
		quizzes = append(quizzes, quiz)
	}

	// Respond with the quizzes as JSON
	c.JSON(http.StatusOK, quizzes)
}

func assignQuizToUsers(c *gin.Context) {
	db := connect()

	// Parse the request body to get the quiz ID and user IDs
	var association QuizUserAssociation
	if err := c.ShouldBindJSON(&association); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Check if the quiz ID is valid
	// (You might want to add additional checks here, such as verifying if the quiz exists)
	if association.QuizID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID"})
		return
	}

	// Check if there are any user IDs to assign the quiz
	if len(association.UserIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No user IDs provided"})
		return
	}

	// Prepare the SQL statement to insert quiz-user associations
	stmt, err := db.Prepare("INSERT INTO tblquizzes_tblusers (tblquizzes_id, tblusers_id) VALUES ($1, $2)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare SQL statement"})
		return
	}
	defer stmt.Close()

	// Insert the quiz-user associations into the database
	for _, userID := range association.UserIDs {
		_, err := stmt.Exec(association.QuizID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign quiz to users"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quiz assigned to users successfully"})
	
}



func getQuizCreated(c *gin.Context) {
    db := connect()

    // Get the user ID from the context and convert it to an integer
    userIDRaw, exists := c.Get("id")
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User ID not found in context"})
        return
    }
    userID, ok := userIDRaw.(int)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Initialize a slice to store the quizzes
    var quizzes []Quiz

    // Query the tblquizzes table to fetch the quiz details for the given user ID
    rows, err := db.Query("SELECT tblusers_id, id,created_time, name FROM tblquizzes WHERE tblusers_id = $1", userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quizzes"})
        return
    }
    defer rows.Close()

    // Iterate through the rows and fetch the quiz details
    for rows.Next() {
        var quiz Quiz
        err := rows.Scan(&quiz.UserID,&quiz.ID,&quiz.CreatedTime,&quiz.Name)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quiz details"})
            return
        }

        // Append the quiz to the quizzes slice
        quizzes = append(quizzes, quiz)
    }
	if quizzes==nil{
		id,_:=c.Get("id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch quizzes","admin":id})
		return
	}
    c.JSON(http.StatusOK, quizzes)
}


/*func assignQuiz(c *gin.Context) {
	db:=connect()

	type FormData struct {
		 string `form:"difficulty"`
		Category        string `form:"category"`
		Number     *int    `form:"number,default=-1"`
	}
	//if c.Get("id")


}*/







func verifyToken(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	if tokenString == "" {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}

	// Parse the token and check if it's valid and not expired
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.JSON(401, gin.H{"message": "unexpected signing method"})
			c.Abort()
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return MySigningKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"message": "Invalid or expired token"})
		c.Abort()
		return
	}

	// Access the user information from the token's claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(401, gin.H{"message": "Invalid token claims"})
		c.Abort()
		return
	}

	// Convert the "id" claim to an int
	id, ok := claims["id"].(float64)
	if !ok {
		c.JSON(401, gin.H{"message": "Invalid ID in token"})
		c.Abort()
		return
	}

	// Add the user information to the Gin context
	c.Set("id", int(id))
	c.Set("username", claims["username"].(string))
	c.Set("email", claims["email"].(string))
	c.Set("is_admin", claims["is_admin"].(bool))
	c.Next()
}


func generateUniqueSeed() int64 {
    // Get the current timestamp in nanoseconds
    timestamp := time.Now().UnixNano()

    // Generate a random number between 0 and 999999
    randomNum := rand.Int63n(1000000)

    // Combine the timestamp and random number to create a unique seed
    seed := timestamp + randomNum

    return seed
}






func getTagsFromString(tagNamesStr string, split string) []Tag {
	tagNames := strings.Split(tagNamesStr, split)
	tags := make([]Tag, 0, len(tagNames)) // Use a slice with initial capacity

	for _, tagName := range tagNames {
		tags = append(tags, Tag{Name: tagName})
	}

	return tags
}



func getRandomElements(arr []Question, n int) []Question {
	if n >= len(arr) {
		return arr
	}

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Create a copy of the original array
	copyArr := make([]Question, len(arr))
	copy(copyArr, arr)

	// Shuffle the copy of the array using the Fisher-Yates algorithm
	for i := len(copyArr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		copyArr[i], copyArr[j] = copyArr[j], copyArr[i]
	}

	// Return the first n elements from the shuffled copy
	return copyArr[:n]
}

		
	



func main(){
	router := gin.Default()
	router.GET("/quizassign",verifyToken,getQuizAssign)
	router.GET("/quizcreated",verifyToken,getQuizCreated)
	router.POST("/assign-quiz",verifyToken,assignQuizToUsers)
	//router.GET("/quizcreated",verifyToken,)
    //router.GET("/createquiz",verifyToken,getQuiz)
	router.POST("/savequiz",verifyToken,saveQuiz)
	router.Run("localhost:8000")
	//db:=connect()

}