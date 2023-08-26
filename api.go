package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Address struct {
	Home   string `json:"home"`
	HPhone string
	Office string `json:"office"`
	OPhone string
}

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Dob       string             `json:"dob,omitempty" bson:"dob,omitempty"`
	City      string             `json:"city,omitempty" bson:"city,omitempty"`
	Country   string             `json:"country,omitempty" bson:"country,omitempty"`
	Address   Address            `json:"address,omitempty" bson:"address,omitempty"`
}

func main() {
	fmt.Println("Start the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	router := gin.Default()

	router.POST("/customer", CreateCustomer)
	router.GET("/customer", GetCustomers)
	router.GET("/customer/country/:country", GetCustomerByCountry)
	router.GET("/customer/city/:city", GetCustomerByCity)
	router.GET("/customer/name/:name", GetCustomerByName)
	router.GET("/customers/countries/:countries", GetCustomersFromCountries)

	router.Run(":0786")
}

func CreateCustomer(c *gin.Context) {
	var person Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := client.Database("People").Collection("User Detai")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func GetCustomers(c *gin.Context) {
	collection := client.Database("People").Collection("User Detai")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var people []Person
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}

	c.JSON(http.StatusOK, people)
}

func GetCustomerByCountry(c *gin.Context) {
	countryName := c.Param("country")
	collection := client.Database("People").Collection("User Detai")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	filter := bson.M{"country": countryName}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var people []Person
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}

	c.JSON(http.StatusOK, people)
}
func GetCustomerByCity(c *gin.Context) {
	cityName := c.Param("city")
	collection := client.Database("People").Collection("User Detai")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	filter := bson.M{"city": cityName}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var people []Person
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}

	c.JSON(http.StatusOK, people)
}


func GetCustomerByName(c *gin.Context) {
	name := c.Param("name")
	collection := client.Database("People").Collection("User Detai")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	
	filter := bson.M{
		"$or": []bson.M{
			{"firstname": bson.M{"$regex": name, "$options": "i"}},
			{"lastname": bson.M{"$regex": name, "$options": "i"}},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var people []Person
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, people)
}
func GetCustomersFromCountries(c *gin.Context) {
	countryNames := c.Param("countries")
	countries := strings.Split(countryNames, ",")
	collection := client.Database("People").Collection("User Detai")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	filter := bson.M{"country": bson.M{"$in": countries}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var people []Person
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}

	c.JSON(http.StatusOK, people)
}
