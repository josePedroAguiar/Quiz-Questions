
package main

import (
	"net/http"
	"fmt"
	"math/rand"
	"time"
	"strconv"
    "os"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

var MySigningKey = []byte(os.Getenv ("SECRET_KEY"))// Replace this with your secret key




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

// Example handler that requires authenticated user information
func protectedHandler(c *gin.Context) {
	// Access user information from Gin context
	username, _ := c.Get("username")
	email, _ := c.Get("email")
	isAdmin, _ := c.Get("is_admin")

	c.JSON(200, gin.H{
		"message":  "Protected endpoint",
		"username": username,
		"email":    email,
		"is_admin": isAdmin,
	})
}

func AllQuestions(c *gin.Context){
	db:=connect()	
	c.IndentedJSON(http.StatusOK,getAll(db))
}



func postQuestionsBy(c *gin.Context) {
	db:=connect()
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


	// DEBUG
	c.JSON(http.StatusOK, gin.H{"message": "Data submitted successfully", "name": formData.Difficulty, "category": formData.Category,"number": formData.Number})
	


	if formData.Difficulty == "" && formData.Category == ""{
		arr:=getAll(db)
		if(arr==nil){
			c.JSON(http.StatusOK, gin.H{"message": "Not found any questions"})
			return
		}
		if formData.Number != nil  && len(arr)>*formData.Number {
			resp:=getRandomElements(arr, *formData.Number)
			c.IndentedJSON(http.StatusOK,resp)
			return
		} else {
			c.IndentedJSON(http.StatusOK,arr)
			return
		}


	}else if formData.Difficulty == ""{
		arr:=getByCategory(db, formData.Category)
		if(arr==nil){
			c.JSON(http.StatusOK, gin.H{"message": "Not found any questions"})
			return
		}
		if formData.Number != nil  && len(arr)>*formData.Number {
			resp:=getRandomElements(arr, *formData.Number)
			c.IndentedJSON(http.StatusOK,resp)
			return
		} else {
			c.IndentedJSON(http.StatusOK,arr)
			return
		}


	}else if formData.Category == ""{
		arr:=getByDifficulty(db,formData.Difficulty)
		if(arr==nil){
			c.JSON(http.StatusOK, gin.H{"message": "Not found any questions"})
			return
		}
		if formData.Number != nil  && len(arr)>*formData.Number  {
			resp:=getRandomElements(arr, *formData.Number)
			c.IndentedJSON(http.StatusOK,resp)
			return
		} else {
			c.IndentedJSON(http.StatusOK,arr)
			return
		}
	}else{
		arr:=getByCategoryAndDifficulty(db,formData.Category,formData.Difficulty)
		if(arr==nil){
			c.JSON(http.StatusOK, gin.H{"message": "Not found any questions"})
			return
		}
		if formData.Number != nil  && len(arr)>*formData.Number {
			resp:=getRandomElements(arr, *formData.Number)
			c.IndentedJSON(http.StatusOK,resp)
			return
		} else {
			c.IndentedJSON(http.StatusOK,arr)
			return
		}
	}



}

func getQuestionsBy(c *gin.Context) {
	// Access query parameters using c.Query()
	difficulty := c.Query("difficulty")
	category := c.Query("category")
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
		number=0
	}

	db := connect()

	if difficulty == "" && category == "" {
		arr := getAll(db)
		if(arr==nil){
			c.JSON(http.StatusOK, gin.H{"message": "Not found any questions"})
			return
		}
		if number != 0  && len(arr)>number{
			resp := getRandomElements(arr, number)
			c.JSON(http.StatusOK, resp)
			return
		} else {
			c.JSON(http.StatusOK, arr)
			return
		}
	} else if difficulty == "" {
		arr := getByCategory(db, category)
		if(arr==nil){
			c.JSON(http.StatusOK, gin.H{"message": "Not found any questions"})
			return
		}
		if number != 0  && len(arr)>number{
			resp := getRandomElements(arr, number)
			c.JSON(http.StatusOK, resp)
			return
		} else {
			c.JSON(http.StatusOK, arr)
			return
		}
	} else if category == "" {
		arr := getByDifficulty(db, difficulty)
		if number != 0  && len(arr)>number{
			resp := getRandomElements(arr, number)
			c.JSON(http.StatusOK, resp)
			return
		} else {
			c.JSON(http.StatusOK, arr)
			return
		}
	} else {
		arr := getByCategoryAndDifficulty(db, category, difficulty)
		if(arr==nil){
			c.JSON(http.StatusOK, gin.H{"message": "Not found any questions"})
			return
		}
		if number != 0 && len(arr)>number {
			resp := getRandomElements(arr, number)
			c.JSON(http.StatusOK, resp)
			return
		} else {
			c.JSON(http.StatusOK, arr)
			return
		}
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



	cats:=getTagsFromString(categorys,",")
	db := connect()

	var questions []Question
	for _,category := range cats {
		// Set the query parameters
		if difficulty == "" && category.Name == "" {
			arr := getAll(db)
			if(arr==nil){
				c.JSON(http.StatusOK, gin.H{"message": "Not found any questions"})
				continue
			}
			if number != 0  && len(arr)>number{
				resp := getRandomElements(arr, number)
				questions = append(questions, resp...)
				continue
			} else {
				questions = append(questions, arr...)
				continue
			}
		} else if difficulty == "" {
			arr := getByCategory(db, category.Name)
			if(arr==nil){
				c.JSON(http.StatusOK, gin.H{"message": "Not found any questions"})
				continue
			}
			if number != 0  && len(arr)>number{
				resp := getRandomElements(arr, number)
				questions = append(questions, resp...)
				continue
			} else {
				c.JSON(http.StatusOK, arr)
				continue
			}
		} else if category.Name == "" {
			arr := getByDifficulty(db, difficulty)
			if number != 0  && len(arr)>number{
				resp := getRandomElements(arr, number)
				questions = append(questions, resp...)
				continue
			} else {
				c.JSON(http.StatusOK, arr)
				continue
			}
		} else {
			arr := getByCategoryAndDifficulty(db, category.Name, difficulty)
			if(arr==nil){
				c.JSON(http.StatusOK, gin.H{"message": "Not found any questions"})
				continue
			}
			if number != 0 && len(arr)>number {
				resp := getRandomElements(arr, number)
				questions = append(questions, resp...)
				continue
			} else {
				questions = append(questions, arr...)
				continue
			}
	
		}
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





func main() {
	db:=connect()

	router := gin.Default()
    router.GET("/allquestions",verifyToken,AllQuestions)
	router.POST("/questions",verifyToken,postQuestionsBy)
	router.GET("/questions",verifyToken,getQuestionsBy)
    router.GET("/quiz",verifyToken,getQuiz)
    router.Run("localhost:8080")
	getAll(db)
	getByDifficulty(db,"'' or 1=1")
}
