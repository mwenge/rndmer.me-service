package main

import (
	"encoding/base64"
  "io/ioutil"
	"net/http"
  "log"

	"github.com/gin-gonic/gin"
)

// video represents data about a record video.
type video struct {
	Date  string  `json:"date"`
	Content string  `json:"content"`
}
var previousVideo video

func SetUpRouter() *gin.Engine {
  dat, err := ioutil.ReadFile("klaus.mp4")
  if err != nil {
    log.Fatal(err)
  }
	encodedString := base64.StdEncoding.EncodeToString(dat)
  previousVideo = video{Date: "2022-05-22", Content: encodedString}

	router := gin.Default()
	router.POST("/video", PostVideo)
	return router
}
func main() {
	router := SetUpRouter()
	router.Run("localhost:8080")
}

// postVideo adds an video from JSON received in the request body.
func PostVideo(c *gin.Context) {
	var newVideo video

	// Call BindJSON to bind the received JSON to
	// newVideo.
	if err := c.BindJSON(&newVideo); err != nil {
		return
	}

	// Add the new video to the slice.
	c.IndentedJSON(http.StatusCreated, previousVideo)
	previousVideo = newVideo
}

