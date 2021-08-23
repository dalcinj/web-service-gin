package main

// full tutorial and code at https://golang.org/doc/tutorial/web-service-gin

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocarina/gocsv"
)

type Client struct { // Our example struct, you can use "-" to ignore a field
	Id       string  `csv:"client_id"`
	Name     string  `csv:"client_name"`
	Age      int64   `csv:"client_age"`
	NotUsed  string  `csv:"-"`
	Score    float64 `csv:"-"`
	Viable   bool    `csv:"-"`
	Priority float64 `csv:"-"`
}

func main() {
	router := gin.Default()
	router.POST("upload", uploadFile)

	router.Run("localhost:8080")
}

func uploadFile(context *gin.Context) {

	context.Request.ParseMultipartForm(10 << 20) // max upload 10Mb file size

	file, fileHeader, openErr := context.Request.FormFile("file")
	if openErr != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	defer file.Close()

	fileName := fileHeader.Filename
	if CheckExtention(fileName) == false {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("invalid file extention, please upload one of the following: '%v'", ".csv"),
			"error":   true,
		})
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		// handle error
	}

	clients := []*Client{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = ';'
		return r
	})

	if err := gocsv.UnmarshalBytes(fileBytes, &clients); err != nil {
		panic(err)
	}

	handleData(clients)

	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("'%s' uploaded", fileName),
	})

}

func handleData(clients []*Client) {
	for _, client := range clients {

		fmt.Println("Age", client.Age)
		client.Score = (2.5*2.8)/4 + float64(client.Age)*2.0 + 16.54
		fmt.Println("Score: ", client.Score)

		min, max := getPricing(client.Score, 2.25)
		fmt.Printf("min value is '%f' max value is '%f'", min, max)

		client.Viable = checkViable()

		if client.Viable == true {
			client.Priority = calculatePriority()
		}
	}
}

func getPricing(paidAmount float64, totalAmount float64) (float64, float64) {
	_, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		panic(err)
	}
	// fmt.Println("response, ", response)
	return (paidAmount + totalAmount), (paidAmount * totalAmount)
}

func checkViable() bool {
	return true
}

func calculatePriority() float64 {
	return 1.0
}
