/*
rloot-server

Copyright © 2023 Austin Jackson <vesche@protonmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
		startServer()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// TODO: remove full path here
var DATABASE_PATH string = "/home/vesche/rloot/rloot-server/test/rloot/_meta/rloot.db"

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

type movie struct {
	gorm.Model
	Title         string `json:"title"`
	Year          int    `json:"year"`
	Director      string `json:"director"`
	Cast          string `json:"cast"`
	TrailerURL    string `json:"trailer_url"`
	MediaFileName string `json:"media_file_name"`
	ImageFileName string `json:"image_file_name"`
}

type tv struct {
	gorm.Model
	SeriesTitle   string `json:"series_title"`
	EpisodeTitle  string `json:"episode_title"`
	SeasonNumber  int    `json:"season_number"`
	EpisodeNumber int    `json:"episode_number"`
	Year          int    `json:"year"`
	Director      string `json:"director"`
	Cast          string `json:"cast"`
	MediaFileName string `json:"media_file_name"`
	ImageFileName string `json:"image_file_name"`
}

type book struct {
	gorm.Model
	Title         string `json:"title"`
	Author        string `json:"author"`
	Year          int    `json:"year"`
	Summary       string `json:"summary"`
	MediaFileName string `json:"media_file_name"`
	ImageFileName string `json:"image_file_name"`
}

func getLibraryData(libraryName string) []byte {
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
		db.Find(&movies)
		j, _ := json.Marshal(movies)
		return j
	} else if libraryName == "tv" {
		var shows []tv
		db.Find(&shows)
		j, _ := json.Marshal(shows)
		return j
	} else if libraryName == "books" {
		var books []book
		db.Find(&books)
		j, _ := json.Marshal(books)
		return j
	} else {
		return []byte(`{"error": "Library not found"}`)
	}
}

func listLibrary(c *gin.Context) {
	libraryName := c.Param("libraryName")
	c.Data(http.StatusOK, gin.MIMEJSON, getLibraryData(libraryName))
}

func getMedia(c *gin.Context) {
	libraryName := c.Param("libraryName")
	fileName := c.Param("fileName")
	// TODO: remove full path here
	filePath := fmt.Sprintf("/home/vesche/rloot/rloot-server/test/rloot/%s/%s", libraryName, fileName)
	if _, err := os.Stat(filePath); err != nil {
		fmt.Println("foo")
	}
	c.File(filePath)
}

func startServer() {
	db, err := gorm.Open(sqlite.Open(DATABASE_PATH), &gorm.Config{})
	if err != nil {
		log.Fatal("foo1")
	}

	db.AutoMigrate(&movie{})
	db.AutoMigrate(&tv{})
	db.AutoMigrate(&book{})

	// example entries
	db.Create(&movie{
		Title:         "Gladiator",
		Year:          2000,
		Director:      "Ridley Scott",
		Cast:          "Russell Crowe, Joaquin Phoenix, Connie Nielsen",
		TrailerURL:    "https://www.youtube.com/watch?v=uvbavW31adA",
		MediaFileName: "Gladiator_2000_1080p.mkv",
		ImageFileName: "gladiator_poster.png",
	})
	db.Create(&tv{
		SeriesTitle:   "Lost",
		EpisodeTitle:  "Man of Science, Man of Faith",
		SeasonNumber:  2,
		EpisodeNumber: 1,
		Year:          2005,
		Director:      "Jack Bender",
		Cast:          "Matthew Fox, Terry O'Quinn, Evangeline Lilly",
		MediaFileName: "Lost_s02e01_720p.mkv",
		ImageFileName: "Lost_poster.png",
	})
	db.Create(&book{
		Title:         "Harry Potter and the Sorcerer's Stone",
		Author:        "J.K. Rowling",
		Year:          1998,
		Summary:       "Harry Potter has no idea how famous he is. That's because he's being raised by his miserable aunt and uncle who are terrified Harry will learn that he's really a wizard, just as his parents were. But everything changes when Harry is summoned to attend an infamous school for wizards, and he begins to discover some clues about his illustrious birthright. From the surprising way he is greeted by a lovable giant, to the unique curriculum and colorful faculty at his unusual school, Harry finds himself drawn deep inside a mystical world he never knew existed and closer to his own noble destiny.",
		MediaFileName: "Harry_Potter_1.pdf",
		ImageFileName: "Harry_Potter_1_poster.png",
	})
	// -----------

	router := gin.Default()

	//router.GET("/ping", func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "pong",
	//	})
	//})

	// list library (movies, tv, anime, video, books, pictures)
	router.GET("/list-library/:libraryName", listLibrary)

	// get media file
	router.GET("/get-media/:libraryName/:fileName", getMedia)

	router.Run()
}
