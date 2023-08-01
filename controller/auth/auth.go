package auth

import (
	"net/http"
	"os"
	"time"

	"golang-apiuser/orm"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterModel struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FullName  string `json:"fullname" binding:"required"`
	AccountNo string `json:"accountno" binding:"required"`
}

// Register			godoc
// @Summary 		Register
// @Description 	User Register API
// @Param			parameter body RegisterModel true "Register"
// @Produce			application/json
// @Tags			User API
// @Router 			/register [post]
func Register(c *gin.Context) {
	var post RegisterModel
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userExist orm.User
	orm.Db.Where("username = ?", post.Username).First(&userExist)
	if userExist.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Username already in exist",
		})
		return
	}

	encryptPassword, _ := bcrypt.GenerateFromPassword([]byte(post.Password), 4)
	user := orm.User{
		Username:  post.Username,
		Password:  string(encryptPassword),
		FullName:  post.FullName,
		AccountNo: post.AccountNo,
		Credit:    1000,
	}
	orm.Db.Create(&user)
	if user.ID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "User added",
			"userId":  user.ID,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "User failed",
		})
	}
}

type LoginModel struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var hmacSampleSecret []byte

// Login			godoc
// @Summary 		Login
// @Description 	User Login API
// @Param			parameter body LoginModel true "Login"
// @Produce			application/json
// @Tags			User API
// @Router 			/login [post]
func Login(c *gin.Context) {

	var post LoginModel
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var usernameExist orm.User
	orm.Db.Where("username = ?", post.Username).First(&usernameExist)
	if usernameExist.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Username not found",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(usernameExist.Password), []byte(post.Password))
	if err == nil {
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": usernameExist.ID,
			"exp":    time.Now().Add(time.Minute * 60).Unix(),
		})

		tokenString, _ := token.SignedString(hmacSampleSecret)
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Login Success",
			"token":   tokenString,
		})

	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "error",
			"message": "Login Fail",
		})
	}

}
