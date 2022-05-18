package main

import (
	"encoding/base64"
  "io/ioutil"
	"net/http"
  "log"
  "time"

  "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// video represents data about a record video.
type video struct {
	Date  string  `json:"date"`
	Content string  `json:"content"`
}

var videoQueue = make(chan video, 300)

func SetUpRouter() *gin.Engine {
  dat, err := ioutil.ReadFile("klaus.mp4")
  if err != nil {
    log.Fatal(err)
  }
	encodedString := base64.StdEncoding.EncodeToString(dat)
  initialVideo := video{Date: "2022-05-22", Content: encodedString}
  videoQueue <- initialVideo

	router := gin.Default()
	router.POST("/video", PostVideo)
	return router
}

func main() {
	router := SetUpRouter()
  router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:4001"},
    AllowMethods:     []string{"POST"},
    AllowHeaders:     []string{"Origin", "Content-Type"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
  }))
	router.Run("localhost:8080")
}

// postVideo adds an video from JSON received in the request body.
func PostVideo(c *gin.Context) {

	var newVideo video
	if err := c.BindJSON(&newVideo); err != nil {
    log.Printf("Unable to Bind JSON")
    s, _ := c.GetRawData()
    log.Printf(string(s))
		return
	}
  c.Header("Access-Control-Allow-Origin", "http://localhost:4001")
  c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT")
  c.Header("Access-Control-Allow-Headers", "Content-Type")
  // Pop a video from the top of the queue.
  select {
  case previousVideo := <-videoQueue:
    log.Printf("Creating Response")
    c.JSON(http.StatusCreated, previousVideo)
  default:
    log.Printf("QUEUE IS EMPTY!")
  }
  videoQueue <- newVideo
}

