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

type Status string

const (
	Initiated   Status = "Initiated"
	Ongoing     Status = "Ongoing"
	UserDropped Status = "UserDropped"
	Cancelled   Status = "Cancelled"
)

//STATUS = > INITIATED, ONGOING, USER_DROPPED, CANCELLED

// Create policy
// 1. add one more field in policy represents Beneficiary, SSN
// 2  premium amount should be validated against the min and max amount.
// 3  premium months should be validated against the min and max period.
// 4. update the policy unique identifier to len(policies purchased by the user) + 1
// 5. calculate the emi_amount using frequency * (premium / months)
// 6. update the valid till by calculating the current time + months
// 5. update status to initiated
// 6. update total coverage to => premium * 0.3 =>
// 7. to generate illustration, iterate a for loop till emi months
//	 a. iterate using the freqency
//for each frequency add illustration with following details
//	1. amount - emi_amount
//	2. date to pay - current date
//	3. total_amount = number of installment * emi_amount
//  4. paid_date -
//	4. status - pending
//
//	1. amount - emi_amount
//	2. date to pay - current date + frequency
//	3. total_amount = number of installment * emi_amount
//  4. paid_date -
//	4. status - pending

// pay first emi
// /users/:email/policy/:1/pay [POST]
// should bind request body
// do validate request body
// should bind URI request
// do validate URI request
// identify the emi(1st / 2nd)
// check how many emi user paid
// if it is first then allow user to pay first emi
// if it is subsequent then check current date is within - 10 and + 10 days of date to pay
// else throw error if user is trying to pay before 10 days of date to pay -(please wait till some date)
// user is paying after 10 days of pay -> please pay installment_amount + (0.1 * number of delayed days)
// update the status to paid
// update the paid_date to current date

//premium = 120000 // user provided be should be validated against hte min and max period.
//months - 14     // user provided be should be validated against hte min and max period.
//frequency = > 3 // user provided
//emi_amount = > frequency * (120000 / 14) // calculated field => frequency * (premium / months)
//valid till = > curernt dat + 14 months // calculated field

type Policy struct {
	UniqueIdentifier int         `json:"unique_identifier"`
	ID               string      `json:"id" validate:"required"`
	Name             string      `json:"name" validate:"required,min=4,max=20"`
	ValidTill        time.Time   `json:"time_period" validate:"required"`
	Months           int         `json:"months" validate:"required"`
	EMIAmount        int         `json:"emi_amount" validate:"required"`
	Status           Status      `json:"status" validate:"required"`
	TotalCoverage    int         `json:"total_coverage" validate:"required"`
	Premium          float64     `json:"premium" validate:"required"`
	Frequency        int         `json:"frequency" validate:"required"`
	Beneficiary      Beneficiary `json:"beneficiary" validate:"required"`
}

type Beneficiary struct {
	Name string `json:"name" validate:"required,min=3,max=40"`
	SSN  string `json:"SSN" validate:"required,ssn"`
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
