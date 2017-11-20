package rapidvideo

import (
  "testing"
  "os"
)

func NewTestRapidVideo() (rv Rapidvideo){
  var userId = os.Getenv("RAPID_VIDEO_USER_ID")
  var apiKey = os.Getenv("RAPID_VIDEO_API_KEY")

  return NewRapidVideo(userId, apiKey)
}

func TestCreateRapidVideo(t *testing.T){
  var userId = os.Getenv("RAPID_VIDEO_USER_ID")
  var apiKey = os.Getenv("RAPID_VIDEO_API_KEY")
  var rv Rapidvideo

  if userId == "" {
    t.Error("RAPID_VIDEO_USER_ID is not defined")
  }

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

func TestGetInfo(t *testing.T) {
  var rv Rapidvideo = NewTestRapidVideo()

  v, err := rv.GetInfo("code")
  if err != nil {
    t.Error("Error upload a file", err.Error())
  }

  if v.Code == "" {
    t.Error("Empty video code")
  }
}

func TestUpload(t *testing.T) {
  var rv Rapidvideo = NewTestRapidVideo()

  err := rv.Upload("")
  if err != nil {
    t.Error("Error upload a file", err.Error())
  }
}

func TestRemoteUpload(t *testing.T) {
  var rv Rapidvideo = NewTestRapidVideo()

  ok, err := rv.RemoteUpload("url")
  if err != nil {
    t.Error("Error creating a remote file upload", err.Error())
  }

  if ok != true {
    t.Error("Should create a remote upload")
  }
}

func TestRemoteStatus(t *testing.T) {
  var rv Rapidvideo = NewTestRapidVideo()

  status, err := rv.RemoteStatus("id")
  if err != nil {
    t.Error("Error getting a remote file upload status", err.Error())
  }

  if status.ObjectCode == "" {
    t.Error("Should get an object code")
  }
}
