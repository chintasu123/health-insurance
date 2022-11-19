package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

type User struct {
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

func main() {
	server := gin.New()
	err := server.Run()
	if err != nil {
		log.Println("unable to server :", err)
	}
}
