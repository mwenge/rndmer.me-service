package main

import (
  "io/ioutil"
  "net/http"
  "log"
  "time"

  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/unrolled/secure"
)

// video represents data about a record video.
type video struct {
  Date  string  `json:"date"`
  Content []byte  `json:"content"`
}

var videoQueue = make(chan video, 300)

func SetUpRouter() *gin.Engine {
  dat, err := ioutil.ReadFile("klaus.mp4")
  if err != nil {
    log.Fatal(err)
  }
  initialVideo := video{Date: "2022-05-22", Content: dat}
  videoQueue <- initialVideo

  router := gin.Default()
  router.POST("/video", PostVideo)
  return router
}

func LoadTls() gin.HandlerFunc {
  return func(c *gin.Context) {
    middleware := secure.New(secure.Options{
      SSLRedirect: true,
      SSLHost:     "localhost:8000",
    })
    err := middleware.Process(c.Writer, c.Request)
    if err != nil {
      //If an error occurs, do not continue.
      log.Println(err)
      return
    }
    //Continue processing
    c.Next()
  }
}

func main() {
  router := SetUpRouter()
  router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:4001", "https://localhost:4443", "https://192.168.192.24:4443"},
    AllowMethods:     []string{"POST"},
    AllowHeaders:     []string{"Origin", "Content-Type"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
  }))
  //router.Run("localhost:8080")
  router.Use(LoadTls())
  //Enable port listening
  router.RunTLS(":8080", "server.pem", "key.pem")
}

// postVideo adds an video from JSON received in the request body.
func PostVideo(c *gin.Context) {

  data, err := c.GetRawData()
  if err != nil {
    log.Printf("Can't get data")
    return
  }
  newVideo := video{Date: "2022-05-22", Content: data}
  c.Header("Access-Control-Allow-Origin", "http://localhost:4001")
  c.Header("Access-Control-Allow-Origin", "https://localhost:4443")
  c.Header("Access-Control-Allow-Origin", "https://192.168.192.24:4443")
  c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT")
  c.Header("Access-Control-Allow-Headers", "Content-Type")
  // Pop a video from the top of the queue.
  select {
  case previousVideo := <-videoQueue:
    log.Printf("Creating Response")
    c.Data(http.StatusCreated, "video/mp4", previousVideo.Content)
  default:
    log.Printf("QUEUE IS EMPTY!")
  }
  videoQueue <- newVideo
}

