package rapidvideo

import (
  "testing"
  "os"
)

func NewTestRapidVideo() (rv Rapidvideo){
  var apiKey = os.Getenv("RAPID_VIDEO_API_KEY")

  return NewRapidVideo(apiKey)
}

func TestCreateRapidVideo(t *testing.T){
  var apiKey = os.Getenv("RAPID_VIDEO_API_KEY")
  var rv Rapidvideo

  if apiKey == "" {
    t.Error("RAPID_VIDEO_API_KEY is not defined")
  }

  rv = NewTestRapidVideo()
  if rv == nil {
    t.Error("Error creating Rapidvideo")
  }

  _, ok := rv.(*rapidvideo)
  if ok == false {
    t.Error("Could not switch Rapidvideo interface to rapidvideo struct")
  }
}

func TestSetProxy(t *testing.T){
  var rv Rapidvideo = NewTestRapidVideo()
  var err error

  err = rv.SetProxy("socks5://127.0.0.1:9050")
  if err != nil {
    t.Error("Error setting a proxy:", err)
  }
}
