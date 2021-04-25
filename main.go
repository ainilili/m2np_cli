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

	// API.get("/search?username=" + txt)
	// API.post("update_reply", { id: reply.id, content: txt })
	// API.get("/my_emojis", {})
	// API.post("/delete_post", { id: article.id })
	// API.post("/update_personal_info", {
	// 	description: description,
	// 	color: color,
	// 	display_name: displayName,
	//   })
	// API.post("/reset-password", { email: email, password: password })
	// API.get("/notifications")
	// API.get("/get_replies", { post_id: post_id })
	// API.get("/get_profile", {username: username})
	// API.get(
	// username == null ? "/get_inbox" : "/get_outbox",
	// mode == "append"
	// 	? { less_than_ts: minTS, username: username }
	// 	: mode == "prepend"
	// 	? { more_than_ts: maxTS, username: username }
	// 	: { username: username }

	// API.post("/unfollow", { user_id: user.id })
	// API.post("/follow", { user_id: user.id })
}
