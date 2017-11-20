package rapidvideo

import (
  "bytes"
  "errors"
  "encoding/json"
  "fmt"
  "golang.org/x/net/proxy"
  "io"
  "io/ioutil"
  "mime/multipart"
  "net/http"
  "net/url"
  "os"
)

type response struct {
  Status int `json: "status"`
  Msg string `json: "msg"`
  Result interface{} `json: "result"`
}

type videoMeta struct {
  Code struct {
    Code string `json: "code"`
    Name string `json: "name"`
    Description string `json: "description"`
  } `json: "code"`
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
  RemoteUpload(u string) (ok bool, err error)

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
  resp, err := rv.get(fmt.Sprintf("/objects.php?ac=info&apikey=%v&code=%v", rv.apiKey, code))
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
  var buf bytes.Buffer
  var req *http.Request
  var writer *multipart.Writer
  var ioWriter io.Writer
  var f *os.File

  writer = multipart.NewWriter(&buf)
  f, err = os.Open(path)
  if err != nil {
    return err
  }
  defer f.Close()

  ioWriter, err = writer.CreateFormFile("file", path)
  if err != nil {
    return err
  }

  if _, err = io.Copy(ioWriter, f); err != nil {
    return err
  }

  if ioWriter, err = writer.CreateFormField("user_id"); err != nil {
    return err
  }
  if _, err = ioWriter.Write([]byte(rv.userId)); err != nil {
    return err
  }
  writer.Close()

  req, err = http.NewRequest("POST", "https://upload.rapidvideo.com/upload.rapidvideo.com/upload/index.php", &buf)
  if err != nil {
    return err
  }

  req.Header.Set("Content-Type", writer.FormDataContentType())
  res, err := rv.client.Do(req)
  if err != nil {
    return err
  }

  if res.StatusCode != 200 {
    return errors.New(fmt.Sprintf("Error uploading a file, expecting 200 got %v", res.StatusCode))
  }

  return err
}

// RemoteUpload remote upload a file
func (rv *rapidvideo) RemoteUpload(u string) (ok bool, err error) {
  resp, err := rv.get(fmt.Sprintf("/remote.php?ac=add&user_id=%v&url=%v", rv.userId, u))
  if err != nil {
    return ok, err
  }

  if resp.Status != 200 {
    return ok, errors.New(resp.Msg)
  }

  return true, err
}

// RemoteStatus check the status of the remote upload
func (rv *rapidvideo) RemoteStatus(id string) (status *uploadStatus, err error){
  resp, err := rv.get(fmt.Sprintf("/remote.php?ac=check&user_id=%v&remote_id=%v", rv.userId, id))
  if err != nil {
    return status, err
  }

  if resp.Status != 200 {
    return status, err
  }

  body, err := json.Marshal(resp.Result)
  if err != nil {
    return status, err
  }

  if err = json.Unmarshal(body, &status); err != nil {
    return status, err
  }

  return status, err
}
