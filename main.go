package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var DATABASE_PATH string = "./rloot.db"

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

type movie struct {
	gorm.Model
	Title          string `json:"title"`
	Year           int    `json:"year"`
	MediaFileName  string `json:"media_file_name"`
	PosterFileName string `json:"poster_file_name"`
	Director       string `json:"director"`
	TrailerURL     string `json:"trailer_url"`
}

type tv struct {
	gorm.Model
	SeriesTitle     string `json:"series_title"`
	EpisodeTitle    string `json:"episode_title"`
	SeasonNumber    int    `json:"season_number"`
	EpisodeNumber   int    `json:"episode_number"`
	Year            int    `json:"year"`
	PrimaryFileName string `json:"primary_file_name"`
	PosterFileName  string `json:"poster_file_name"`
	Director        string `json:"director"`
}

func getLibraryData(libraryName string) string {
	db, err := gorm.Open(sqlite.Open(DATABASE_PATH), &gorm.Config{})
	if err != nil {
		log.WithFields(log.Fields{
			"libraryName": libraryName,
		}).Fatal("Failed to connect to database")
	}
	if err != nil {
		log.WithFields(log.Fields{
			"libraryName": libraryName,
		}).Fatal("Could not open sqlite table")
	}
	if libraryName == "movies" {
		var movies []movie
		result := db.Find(&movies)
		fmt.Println(movies)
		fmt.Println(result)
		j, _ := json.Marshal(movies)
		return string(j)
	} else if libraryName == "tv" {
		var shows []tv
		result := db.Find(&shows)
		fmt.Println(result)
		j, _ := json.Marshal(shows)
		return string(j)
	} else {
		return "{\"error\": \"foo\"}"
	}
}

func listLibrary(c *gin.Context) {
	libraryName := c.Param("libraryName")
	c.IndentedJSON(http.StatusOK, getLibraryData(libraryName))
}

func getMedia(c *gin.Context) {
	libraryName := c.Param("libraryName")
	fileName := c.Param("fileName")
	// TODO: remove test here
	filePath := fmt.Sprintf("test/%s/%s", libraryName, fileName)
	if _, err := os.Stat(filePath); err != nil {
		fmt.Println("foo")
	}
	c.File(filePath)
}

func main() {
	db, err := gorm.Open(sqlite.Open(DATABASE_PATH), &gorm.Config{})
	if err != nil {
		log.Fatal("foo1")
	}

	db.AutoMigrate(&movie{})
	db.AutoMigrate(&tv{})

	// example entries
	db.Create(&movie{
		Title:          "Gladiator",
		Year:           2000,
		MediaFileName:  "Gladiator_2000_1080p.mkv",
		PosterFileName: "gladiator_poster.png",
		Director:       "Ridley Scott",
		TrailerURL:     "https://www.youtube.com/watch?v=uvbavW31adA",
	})
	db.Create(&tv{
		SeriesTitle:     "Lost",
		EpisodeTitle:    "Man of Science, Man of Faith",
		SeasonNumber:    2,
		EpisodeNumber:   1,
		Year:            2005,
		PrimaryFileName: "Lost_s02e01_720p.mkv",
		PosterFileName:  "Lost_poster.png",
		Director:        "Jack Bender",
	})

	// -----------

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// list library (movies, tv, anime, video, books, pictures)
	router.GET("/list-library/:libraryName", listLibrary)

	// get media file
	// fs structure
	// opt/
	//     rloot/
	//         movies/
	//             <media files>
	//             img/
	//                 <poster files>
	//         tv/
	//             <media files>
	//             img/
	//                 <poster files>
	router.GET("/get-media/:libraryName/:fileName", getMedia)

	router.Run()
}
