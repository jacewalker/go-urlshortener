package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ShortUrl struct {
	ShortString string
	WebAddress  string
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "URL Shortener",
		})
	})

	router.POST("/shorten", func(c *gin.Context) {
		webAddress := c.PostForm("url")
		workingUrl := ShortUrl{ShortString: generateRandomString(), WebAddress: formatUrl(webAddress)}

		fmt.Println("Short String: ", workingUrl.ShortString, ", Web Address: ", workingUrl.WebAddress)
		saveToDatabase(workingUrl)

		c.HTML(http.StatusOK, "shorten.tmpl", gin.H{
			"title": "URL Shortener",
		})
	})

	router.Run(":8080")
}

func formatUrl(webAddress string) string {
	u, err := url.Parse(webAddress)

	if err != nil {
		log.Fatal(err)
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		u.Scheme = "https"
	}

	return u.String()
}

func generateRandomString() string {
	// Generates a random string of length 4
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 4)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func saveToDatabase(newUrl ShortUrl) {
	dsn := fmt.Sprintf("host=localhost port=5432 sslmode=disable TimeZone=Australia/Melbourne")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Need to install Postgres container and configure connection details to it

	fmt.Println("DB: ", db)
	fmt.Println("Error: ", err)

	createDatabaseCommand := fmt.Sprintf("CREATE DATABASE URLS")
	db.Exec(createDatabaseCommand)

	if err != nil {
		log.Fatal(err)
	}
	db.Create(&ShortUrl{
		ShortString: newUrl.ShortString,
		WebAddress:  newUrl.WebAddress,
	})
}

// func lookupFromDatabase(shortString string) string {
// 	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=GoUrlShort port=5432 sslmode=disable TimeZone=Australia/Melbourne"), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var shortUrl string
// 	db.Where("short_string = ?", shortString).First(&shortUrl)
// 	return shortUrl
// }
