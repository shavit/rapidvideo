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

type Rapidvideo interface {
  // get make a GET request to the api
  get(endpoint string) (resp *response, err error)

  // SetProxy set a proxy URL
  SetProxy(u string) (err error)
}

type rapidvideo struct {
  apiKey string
  client *http.Client
}

// Craete a new rapidvideo
func NewRapidVideo(apiKey string) Rapidvideo {
  return &rapidvideo{
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
