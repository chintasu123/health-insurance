package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type User struct {
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Gender        string    `json:"gender"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	Age           int       `json:"age"`
	MartialStatus string    `json:"martial_status"`
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

type AvailablePolicy struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	MinAmount     float64 `json:"min_amount"`
	MaxAmount     float64 `json:"max_amount"`
	MinTimePeriod int     `json:"min_time_period"`
	MaxTimePeriod int     `json:"max_time_period"`
}

var (
	users = make(map[string]User)
)

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

	err := server.Run()
	if err != nil {
		log.Println("unable to server :", err)
	}
}
