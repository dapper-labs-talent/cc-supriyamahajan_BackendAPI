package controllers

import (
	"cc-supriyamahajan_BackendAPI/auth"
	"cc-supriyamahajan_BackendAPI/db"
	"cc-supriyamahajan_BackendAPI/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Get all users.
func GetUsers(context *gin.Context) {
	var users []models.User
	err := db.DB.Model(&users).Select()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	var userResponse []models.UserResponse
	for _, u := range users {
		uResponse := models.UserResponse{FirstName: u.FirstName, LastName: u.LastName, Email: u.Email}
		userResponse = append(userResponse, uResponse)
	}
	context.JSON(http.StatusOK, gin.H{"users": userResponse})
}

// SignUp for the new user.
func SignUp(context *gin.Context) {
	var userInfo models.UserInfo
	var user models.User
	if err := context.ShouldBindJSON(&userInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	}

	if err := user.HashPassword(userInfo.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	_, err := db.DB.Model(&models.User{FirstName: userInfo.FirstName,
		LastName: userInfo.LastName,
		Email:    userInfo.Email,
		Password: user.Password,
	}).Insert()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	token, err := auth.GenerateJWT(userInfo.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	context.JSON(http.StatusCreated, gin.H{"token": token})
}

// Login for the registered user.
func Login(context *gin.Context) {
	var userLogin models.UserLogin

	if err := context.ShouldBindJSON(&userLogin); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	var user models.User
	err := db.DB.Model(&user).Where("email = ?", userLogin.Email).Select()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Invalid Credentials"})
		context.Abort()
		return
	}

	token, err := auth.GenerateJWT(user.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	context.JSON(http.StatusOK, gin.H{"token": token})
}

// Update current user first_name and last_name.
func UpdateUser(context *gin.Context) {
	var userName models.UserName

	if err := context.ShouldBindJSON(&userName); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	email, err := getCurrentUserEmail(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	user := models.User{
		FirstName: userName.FirstName,
		LastName:  userName.LastName,
	}

	_, err = db.DB.Model(&user).Column("first_name", "last_name").Where("email = ?0", email).Update()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{})

}

// Get currentUser Email.
func getCurrentUserEmail(context *gin.Context) (string, error) {
	tokenString := context.GetHeader("x-authentication-token")
	claims, err := auth.ParseClaim(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Email, nil
}
