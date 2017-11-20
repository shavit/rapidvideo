package rapidvideo

import (
  "errors"
  "encoding/json"
  "fmt"
  "golang.org/x/net/proxy"
  "io/ioutil"
  "net/http"
  "net/url"
)

type response struct {
  Status int `json: "status"`
  Msg string `json: "msg"`
  Result interface{} `json: "result"`
}

type videoMeta struct {
  Code string `json: "code"`
  Name string `json: "name"`
  Description string `json: "description"`
}

type uploadStatus struct {
  TotalFilesize string `json: "total_filesize"`
  TransferFilesize string `json: "transfer_filesize"`
  Progress string `json: "progress"`
  Done string `json: "done"`
  ObjectCode string `json: "object_code"`
}

type Rapidvideo interface {
  // get make a GET request to the api
  get(endpoint string) (resp *response, err error)

  // SetProxy set a proxy URL
  SetProxy(u string) (err error)

  // GetInfo check if a video is online
  GetInfo(code string) (v *videoMeta, err error)

  // Upload upload a video file and receive the video id
  Upload(path string) (err error)

  // RemoteUpload remote upload a file
  RemoteUpload(url string) (ok bool, err error)

  // RemoteStatus check the status of the remote upload
  RemoteStatus(id string) (status *uploadStatus, err error)
}

type rapidvideo struct {
  userId string
  apiKey string
  client *http.Client
}

// Craete a new rapidvideo
func NewRapidVideo(userId, apiKey string) Rapidvideo {
  return &rapidvideo{
    userId: userId,
    apiKey: apiKey,
    client: new(http.Client),
  }
}

// get make a GET request to the api
func (rv *rapidvideo) get(endpoint string) (resp *response, err error) {
  var res *http.Response
  var body []byte
  resp = new(response)

  res, err = rv.client.Get("https://api.rapidvideo.com/v1" + endpoint)
  if err != nil {
    return resp, err
  }
  defer res.Body.Close()

  if res.StatusCode != 200 {
    return resp, errors.New(fmt.Sprintf("Service unavailable, to status code %d", res.StatusCode))
  }

  body, err = ioutil.ReadAll(res.Body)
  if err != nil {
    return resp, err
  }
  err = json.Unmarshal(body, &resp)

  return resp, err
}

// SetProxy set a proxy URL
func (rv *rapidvideo) SetProxy(u string) (err error) {
  _url, err := url.Parse(u)
  if err != nil {
    return err
  }

  dialer, err := proxy.FromURL(_url, proxy.Direct)
  if err != nil {
    return err
  }

  rv.client.Transport = &http.Transport{Dial: dialer.Dial}

  return err
}

// GetInfo check if a video is online
func (rv *rapidvideo) GetInfo(code string) (v *videoMeta, err error) {
  resp, err := rv.get(fmt.Sprintf("//objects.php?ac=info&apikey=%v&code=%v", rv.apiKey, code))
  if err != nil {
    return v, err
  }

  if resp.Status != 200 {
    return v, errors.New(resp.Msg)
  }

  body, err := json.Marshal(resp.Result)
  if err != nil {
    return v, err
  }

  if err = json.Unmarshal(body, &v); err != nil {
    return v, err
  }

  return v, err
}

// Upload upload a video file and receive the video id
func (rv *rapidvideo) Upload(path string) (err error) {
  return err
}

// RemoteUpload remote upload a file
func (rv *rapidvideo) RemoteUpload(url string) (ok bool, err error) {
  return ok, err
}

// RemoteStatus check the status of the remote upload
func (rv *rapidvideo) RemoteStatus(id string) (status *uploadStatus, err error){
  status = new(uploadStatus)
  
  return status, err
}
