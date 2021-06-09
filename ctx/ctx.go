package ctx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"m2np_cli/model"
	"net/http"
	"time"
)

var httpClient = &http.Client{}

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{
		Transport: tr,
		Timeout:   time.Duration(3) * time.Second,
	}
}

type M2npContext struct {
	User string
	Token string
}

type mss map[string]interface{}

func (x mss) ToJson() *bytes.Buffer {
	data, err := json.Marshal(x)
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewBuffer(data)
}

func (c *M2npContext) Login(user, pwd string) error{
	resp, err := http.Post("https://m2np.com/api/login", "application/json", mss{"email": user, "password": pwd}.ToJson())
	if err != nil{
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var res = map[string]string{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	if res["msg"] != "ok" {
		return errors.New(res["msg"])
	}
	c.Token = res["token"]
	c.User = user
	return nil
}

func (c *M2npContext) GetInBox() (*model.InBoxResp, error){
	request, _ := http.NewRequest("GET", "https://m2np.com/api/get_inbox", nil)
	request.Header.Set("m2np-token", c.Token)
	resp, err := httpClient.Do(request)
	if err != nil{
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var res = &model.InBoxResp{}
	_ = json.NewDecoder(resp.Body).Decode(res)
	return res, nil
}

func (c *M2npContext) PostPost(content string) error{
	request, _ := http.NewRequest("POST", "https://m2np.com/api/post_post", mss{"content": content}.ToJson())
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("m2np-token", c.Token)
	resp, err := httpClient.Do(request)
	if err != nil{
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var res = map[string]string{}
	_ = json.NewDecoder(resp.Body).Decode(&res)
	if res["msg"] != "ok" {
		return errors.New(res["msg"])
	}
	return nil
}