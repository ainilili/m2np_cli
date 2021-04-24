package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type mss map[string]interface{}

func (x mss) ToJson() *bytes.Buffer {
	json_data, err := json.Marshal(x)

	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewBuffer(json_data)
}
func toString(x interface{}) string {
	return fmt.Sprintf("%v", x)
}
func GetContent(method string, url string, msg mss, token string) mss {
	req, _ := http.NewRequest(method, url, msg.ToJson())
	req.Header.Set("m2np-token", token)
	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	return res
}
func main() {
	resp, err := http.Post("https://m2np.com/api/login", "application/json", mss{"email": "abbychau@gmail.com", "password": ""}.ToJson())

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res["token"])
	token := toString(res["token"])

	r := GetContent("GET", "https://m2np.com/api/get_inbox", mss{"username": "abbychau"}, token) //時間線
	fmt.Println(r)
	fmt.Println("----")
	r = GetContent("GET", "https://m2np.com/api/get_outbox", mss{"username": "abbychau"}, token) //username發的文
	fmt.Println(r)
	fmt.Println("----")
}
