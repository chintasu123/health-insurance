package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"net/mail"
	"time"
)

var (
	users = make(map[string]User)
)

var fieldValidator *validator.Validate

type User struct {
	Email         string    `json:"email" validate:"required,email"`
	FirstName     string    `json:"first_name" validate:"required,alpha"`
	LastName      string    `json:"last_name" validate:"required,alpha"`
	Gender        string    `json:"gender" validate:"required,oneof=M F"`
	DateOfBirth   time.Time `json:"date_of_birth" validate:"required"`
	Age           int       `json:"age" validate:"required,min=10,max=60,validateDOB"`
	MartialStatus string    `json:"martial_status" validate:"required,oneof=single married divorce widow"`
	Policy        []Policy  `json:"policy"`
	Address       []Address `json:"address"`
	Premium       Premium   `json:"premium"`
}

type Policy struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Amount        float64   `json:"amount"`
	TimePeriod    time.Time `json:"time_period"`
	EMI           int       `json:"emi"`
	Status        string    `json:"status"`
	TotalCoverage int       `json:"total_coverage"`
}

type Premium struct {
	PaymentType string `json:"payment_type"`
}

type Address struct {
	Line1      string `json:"line_1"`
	Line2      string `json:"line_2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode int32  `json:"postal_code"`
}

type responseMessage struct {
	Message string `json:"message"`
}

func init() {
	fieldValidator = validator.New()
	_ = fieldValidator.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		email := fl.Field().Interface().(string)
		_, err := mail.ParseAddress(email)
		return err == nil
	})
	_ = fieldValidator.RegisterValidation("validateDOB", func(fl validator.FieldLevel) bool {
		age := fl.Field().Interface().(int)
		user := fl.Parent().Interface().(User)
		since := time.Since(user.DateOfBirth)
		days := since.Hours() / 24
		year := int(days / 365)
		return year == age
	})
}

func main() {
	server := gin.New()
	server.POST("/users", func(context *gin.Context) {
		var user User
		err := context.ShouldBindJSON(&user)
		if err != nil {
			log.Println("unable to parse the payload ", err)
			context.JSON(400, map[string]string{
				"message": err.Error(),
			})
			return
		}
		err = fieldValidator.Struct(user)
		if err != nil {
			context.JSON(400, map[string]string{
				"message": err.Error(),
			})
			return
		}

		// user is already existed
		_, isExists := users[user.Email]
		if isExists {
			context.JSON(http.StatusConflict, map[string]string{
				"message": "user already exist",
			})
			return
		}

		// create the user
		users[user.Email] = user

		// return success response
		context.JSON(http.StatusOK, map[string]string{
			"message": "user created",
		})
	})
	server.GET("/users/:email", func(context *gin.Context) {
		email := context.Param("email")
		user, isExists := users[email]
		if !isExists {
			context.JSON(http.StatusNotFound, map[string]string{
				"message": "user not found",
			})
			return
		}
		context.JSON(http.StatusOK, user)
	})
	err := server.Run()
	if err != nil {
		log.Println("unable to server :", err)
	}
}
