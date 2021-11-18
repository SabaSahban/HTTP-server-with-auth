package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	Email string
	jwt.StandardClaims
}

//a struct to rep user account
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

//Validate incoming user details...
func (account *Account) Validate(c *gin.Context) {

	if !strings.Contains(account.Email, "@") {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "email must contain @",
		})
		return
	}

	if len(account.Password) < 6 {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "password is too short",
		})
		return
	}
	//Email must be unique
	temp := &Account{}

	//check for errors and duplicate emails
	err := GetDB().Table("accounts").Where("email = ?", account.Email).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "connection error, try again",
		})
		return
	}

	if temp.Email != "" {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "email address in use",
		})
		return
	}
	account.Create(c)

	c.JSON(200, gin.H{
		"error":   false,
		"message": "requirement passed",
	})
	return
}

func (account *Account) Create(c *gin.Context) {

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		c.JSON(422, gin.H{
			"error":   true,
			"message": "connection error, try again",
		})
	}

	//Create new JWT token for the newly registered account
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString
	account.Password = "" //delete password
	c.JSON(200, gin.H{
		"error":   false,
		"message": "account has been created",
	})
	c.Header("Postman-Token",  account.Token)
	c.JSON(200, account)

}

func Login(email, password string, c *gin.Context) {

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(420, gin.H{
				"error":   true,
				"message": "email address not found",
			})
			return
		}
		c.JSON(420, gin.H{
			"error":   true,
			"message": "connection error, try again",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
			c.JSON(420, gin.H{
				"error":   true,
				"message": "Invalid login credentials. Please try again",
			})
		return
		}
	//Worked! Logged In
	account.Password = ""
	//Create JWT token
	tk := &Token{UserId: account.ID, Email: account.Email}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString //Store the token in the response
	c.JSON(200, gin.H{
		"error":   false,
		"message": "you are logged in",
	})
}

func GetUser(u uint) *Account {
	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { //User not found!
		return nil
	}
	acc.Password = ""
	return acc
}
