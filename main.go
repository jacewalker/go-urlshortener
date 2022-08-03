package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

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
		duplicateCatch, success := lookupWebAddressFromDatabase(workingUrl.WebAddress)

		if success {
			workingUrl = duplicateCatch
		} else {
			saveToDatabase(workingUrl)
		}

		c.HTML(http.StatusOK, "shorten.tmpl", gin.H{
			"title":  "URL Shortener",
			"shorty": fmt.Sprintf("http://127.0.0.1:8080/" + workingUrl.ShortString),
		})
	})

	router.GET("/:shortString", func(c *gin.Context) {
		short := c.Param("shortString")
		if short != "" {
			lookup, success := lookupShortStringFromDatabase(short)

			if success {
				c.Redirect(http.StatusTemporaryRedirect, lookup.WebAddress)
			} else {
				c.HTML(http.StatusOK, "index.tmpl", gin.H{
					"title": "URL Shortener",
				})
			}
		}
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
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:4]
}

func saveToDatabase(newUrl ShortUrl) {
	dsn := fmt.Sprintf("host=172.17.0.3 port=5432 dbname=short_urls user=postgres password=postgres sslmode=disable TimeZone=Australia/Melbourne")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database.")
	}

	db.AutoMigrate(&ShortUrl{})

	// Check if the web address already exists in the database
	var existingUrl ShortUrl
	db.Where("web_address = ?", newUrl.WebAddress).First(&existingUrl)
	if existingUrl.ShortString != "" {
		fmt.Println(newUrl.WebAddress, "already exists in the database.")
	} else {
		db.Create(&ShortUrl{ShortString: newUrl.ShortString, WebAddress: newUrl.WebAddress})
		log.Println("Saved:", newUrl.WebAddress, "to the database.")
	}

}

func lookupShortStringFromDatabase(shortString string) (ShortUrl, bool) {
	dsn := fmt.Sprintf("host=172.17.0.3 port=5432 dbname=short_urls user=postgres password=postgres sslmode=disable TimeZone=Australia/Melbourne")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database.")
	}

	short := new(ShortUrl)
	db.Where("short_string = ?", shortString).First(&short)

	log.Println("Retrieved Short Address:", short.ShortString)
	log.Println("Retrieved Web Address:", short.WebAddress)

	if short.WebAddress == "" {
		return *short, false
	} else {
		return *short, true
	}
}

func lookupWebAddressFromDatabase(webAddress string) (ShortUrl, bool) {
	dsn := fmt.Sprintf("host=172.17.0.3 port=5432 dbname=short_urls user=postgres password=postgres sslmode=disable TimeZone=Australia/Melbourne")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database.")
	}

	short := new(ShortUrl)
	db.Where("web_address = ?", webAddress).First(&short)

	log.Println("Retrieved Short Address:", short.ShortString)
	log.Println("Retrieved Web Address:", short.WebAddress)

	if short.WebAddress == "" {
		return *short, false
	} else {
		return *short, true
	}
}
