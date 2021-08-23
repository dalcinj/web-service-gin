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

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Strain", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry", Price: 17.99},
	{ID: "3", Title: "Sarah", Artist: "Vaughan", Price: 39.99},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlgums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.POST("upload", uploadFile)

	router.Run("localhost:8080")
}

func getAlgums(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(context *gin.Context) {
	var newAlbum album

	if err := context.BindJSON(&newAlbum); err != nil {
		return
	}

	albums = append(albums, newAlbum)
	context.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumByID(context *gin.Context) {
	id := context.Param("id")

	for _, album := range albums {
		if album.ID == id {
			context.IndentedJSON(http.StatusOK, album)
			return
		}
	}

	context.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// type CustomData struct {
// 	id   string
// 	name string
// }

type Client struct { // Our example struct, you can use "-" to ignore a field
	Id      string  `csv:"client_id"`
	Name    string  `csv:"client_name"`
	Age     int64   `csv:"client_age"`
	NotUsed string  `csv:"-"`
	score   float64 `csv:"-"`
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

	// records, readErr := readData(file)

	// if readErr != nil {
	// 	log.Fatal(readErr)
	// 	context.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": readErr.Error(),
	// 		"error":   true,
	// 	})
	// 	return
	// }

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
	for _, client := range clients {
		fmt.Println("Age, ", client.Age)
		client.score = (2.5*2.8)/4 + float64(client.Age)*2.0 + 16.54
		fmt.Println("score, ", client.score)
	}

	// for _, record := range records {
	// 	data := CustomData{
	// 		id:   record[0],
	// 		name: record[1],
	// 	}
	// 	fmt.Printf("my name is %v\n", data.name)
	// }

	context.IndentedJSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("'%s' uploaded", fileName),
	})

	// context.String(http.StatusOK, "ok")

}

// https://medium.com/wesionary-team/easy-working-with-csv-in-golang-using-gocsv-package-9c8424728bbe
// func readData(file multipart.File) ([][]string, error) {

// 	reader := csv.NewReader(file)

// 	// skip first line
// 	if _, err := reader.Read(); err != nil {
// 		return [][]string{}, err
// 	}

// 	reader.Comma = ';'
// 	// reader.Comment = '#'

// 	records, err := reader.ReadAll()

// 	if err != nil {
// 		return [][]string{}, err
// 	}

// 	return records, nil
// }
