package main

import (
  "bytes"
	"encoding/base64"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/stretchr/testify/assert"
)
func EncodeVideo(v string, t *testing.T) string {
  dat, err := ioutil.ReadFile(v)
  if err != nil {
    t.Fatal(err)
  }
	encodedString := base64.StdEncoding.EncodeToString(dat)
  return encodedString
}
func TestPostVideo(t *testing.T) {
	deluge := EncodeVideo("deluge.mp4", t)
  klaus := EncodeVideo("klaus.mp4", t)

  // Test the first call.
  var jsonData = []byte(`{
    "date": "2022-05-22",
    "content": "` + deluge + `"
  }`)

  var expected = []byte(`{
    "date": "2022-05-22",
    "content": "` + klaus + `"
}`)

  router := SetUpRouter()

  w := httptest.NewRecorder()
  req, err := http.NewRequest("POST", "/video", bytes.NewBuffer(jsonData))
  if err != nil {
    t.Fatal(err)
  }

  router.ServeHTTP(w, req)

  assert.Equal(t, 201, w.Code)
  assert.Equal(t, string(expected), w.Body.String())

  // Test the second call, giving us the previously submitted video.
  jsonData = []byte(`{
    "date": "2022-05-22",
    "content": "` + klaus + `"
  }`)

  expected = []byte(`{
    "date": "2022-05-22",
    "content": "` + deluge + `"
}`)

  req, err = http.NewRequest("POST", "/video", bytes.NewBuffer(jsonData))
  if err != nil {
    t.Fatal(err)
  }

  w = httptest.NewRecorder()
  router.ServeHTTP(w, req)

  assert.Equal(t, 201, w.Code)
  assert.Equal(t, string(expected), w.Body.String())
}
