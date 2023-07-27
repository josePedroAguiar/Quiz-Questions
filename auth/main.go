
package main

import (
	"fmt"
	"os"
	_ "github.com/lib/pq"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	jwt "github.com/dgrijalva/jwt-go"
)

var MySigningKey = []byte(os.Getenv ("SECRET_KEY"))// Replace this with your secret key
//var MySigningKey = "ola"

func GetToken(username, email string, isAdmin bool ,id int )string{
	token := jwt.New(jwt.SigningMethodHS256)

	// Create the claims (payload) for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["username"] = username
	claims["email"] = email
	claims["is_admin"] = isAdmin
	//claims["exp"] = time.Now().Add(time.Second*10).Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time (1 day)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(MySigningKey))
	if err != nil {
		return ""
	}

	return tokenString
}


func Login(c *gin.Context) {

	db:=connect()
    type FormData struct {
        Password  string `form:"password" binding:"required"`
        Email string `form:"email" binding:"required,email"`
    }

    var formData FormData

    // Bind the POST data to the struct
    if err := c.ShouldBind(&formData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if the required fields are empty
    if formData.Password == "" || formData.Email == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email and Password are required fields"})
        return
    }

	//user:=getUserByEmailAndPassword(db ,"user@example.com","securepassword")
	user:=getUserByEmailAndPassword(db ,formData.Email,formData.Password)
	
	if user!=nil{
		token := GetToken(user.Username,user.Email,user.IsAdmin,user.ID)
		c.IndentedJSON(http.StatusOK, gin.H{"token": token})
	}else{
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Email or Password are invalid "})
	}
	
    c.JSON(http.StatusOK, gin.H{"message": "Data submitted successfully", "email": formData.Email, "password": formData.Password})
	return
}
func SignUp(c *gin.Context) {

	db:=connect()

    type FormData struct {
        Password  string `form:"password" binding:"required"`
        Email string `form:"email" binding:"required,email"`
		Username string `form:"username" binding:"required"`
    }

    var formData FormData

    // Bind the POST data to the struct
    if err := c.ShouldBind(&formData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if the required fields are empty
    if formData.Password == "" || formData.Email == "" || formData.Username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Email,Username and Password are required fields"})
        return
    }

	//user:=getUserByEmailAndPassword(db ,"user@example.com","securepassword")
	newUser := User{
		Username: formData.Username,
		Email:   formData.Email,
		Password: formData.Password , // Replace this with the actual password input from the user
		IsAdmin:  false,            // Set this to true if the user is an admin
	}

    c.JSON(http.StatusOK, gin.H{"message": AddUser(db,newUser)})
	return
}

func UpdateRoll(c *gin.Context) {

	db:=connect()

	type FormData struct {
        Email string `form:"email" binding:"required,email"`
		IsAdmin string `form:"isadmin" binding:"required"`
    }

	var formData FormData

    if err := c.ShouldBind(&formData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	isAdmin, _:=c.Get("is_admin")
	if isAdmin==true{
		c.JSON(http.StatusBadRequest, gin.H{"message": changeRoll(db,formData.Email,formData.IsAdmin)})
		
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"message": "You can't change Users Rolls","admin":isAdmin})
	}



}

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

		
	



func main(){
	router := gin.Default()
    router.POST("/",Login)
	router.POST("login",Login)
	router.POST("/signup",SignUp)
	router.POST("/updateroll",verifyToken,UpdateRoll)
	router.Run("localhost:8888")
	//db:=connect()
	//getUser(db)
	//fmt.Println(getUserByEmailAndPassword(db ,"user@example.com","securepassword").Username)
}