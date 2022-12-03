package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"net/mail"
	"time"
)

// /users/suvani@gmail.com/policy [POST]
//  should bind uri
// 	do field level validations for uri struct
//  should bind body
//  do field level validations for json struct
//  check if given email is exists or not in the users map
//  if not throw not found error
//  check if the policy is exist in the system
//  add the policy from request to the user list of policy
//  update the user back to the map
//  return 201 policy created message

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
	Policies      []Policy  `json:"policies"`
	Address       []Address `json:"address"`
	Premium       Premium   `json:"premium"`
}

type Policy struct {
	ID            string    `json:"id" validate:"required"`
	Name          string    `json:"name" validate:"required,min=4,max=20"`
	Amount        float64   `json:"amount" validate:"required"`
	TimePeriod    time.Time `json:"time_period" validate:"required"`
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
type EmailURI struct {
	Email string `uri:"email" validate:"required,email"`
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

type AvailablePolicy struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	MinAmount     float64 `json:"min_amount"`
	MaxAmount     float64 `json:"max_amount"`
	MinTimePeriod int     `json:"min_time_period"`
	MaxTimePeriod int     `json:"max_time_period"`
}

type Plan struct {
	Name     string            `json:"name"`
	No       int               `json:"no"`
	Policies []AvailablePolicy `json:"policies"`
}

var (
	plans = []Plan{
		{
			Name: "Super Health Plan",
			No:   1234,
			Policies: []AvailablePolicy{
				{
					ID:            "2345",
					Name:          "whole body",
					MinAmount:     100000.00,
					MaxAmount:     250000.00,
					MinTimePeriod: 12,
					MaxTimePeriod: 24,
				},
			},
		},
		{
			Name: "Classic Health Plan",
			No:   4567,
			Policies: []AvailablePolicy{
				{
					ID:            "7892",
					Name:          "Policy for Eyes",
					MinAmount:     110000.00,
					MaxAmount:     260000.00,
					MinTimePeriod: 10,
					MaxTimePeriod: 23,
				},
			},
		},
	}
)

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

	server.GET("/users/all-plans", func(context *gin.Context) {
		context.JSON(http.StatusOK, plans)
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

	server.POST("/users/:email/policy", func(context *gin.Context) {
		var policy Policy
		var emailURI EmailURI
		err := context.ShouldBindJSON(&policy)
		if err != nil {
			log.Println("unable to parse the payload ", err)
			context.JSON(400, map[string]string{
				"message": err.Error(),
			})
			return
		}
		err = fieldValidator.Struct(policy)
		if err != nil {
			context.JSON(400, map[string]string{
				"message": err.Error(),
			})
			return
		}
		err = context.ShouldBindUri(&emailURI)
		if err != nil {
			log.Println("unable to parse the payload ", err)
			context.JSON(400, map[string]string{
				"message": err.Error(),
			})
			return
		}
		err = fieldValidator.Struct(emailURI)
		if err != nil {
			context.JSON(400, map[string]string{
				"message": err.Error(),
			})
			return
		}
		// user is already existed
		user, isExists := users[emailURI.Email]
		if !isExists {
			context.JSON(http.StatusNotFound, map[string]string{
				"message": "user not found",
			})
			return
		}

		var isPlanPresent bool
		for _, plan := range plans {
			for _, availablePolicy := range plan.Policies {
				if availablePolicy.ID == policy.ID {
					isPlanPresent = true
					break
				}
			}
		}
		if !isPlanPresent {
			log.Println("Policy not found")
			context.JSON(http.StatusNotFound, map[string]string{
				"message": "policy not found",
			})
			return
		}

		user.Policies = append(user.Policies, policy)

		// update the user
		users[user.Email] = user

		// return success response
		context.JSON(http.StatusOK, map[string]string{
			"message": "policy created",
		})
	})
	err := server.Run()
	if err != nil {
		log.Println("unable to server :", err)
	}
}
