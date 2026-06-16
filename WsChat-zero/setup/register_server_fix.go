package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginResp struct {
	Code     int    `json:"code"`
	Token    string `json:"token"`
	UserId   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
}

type CommonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func request(base, path string, body map[string]any) (map[string]any, error) {
	data, _ := json.Marshal(body)
	resp, err := http.Post(base+path, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var result map[string]any
	json.Unmarshal(raw, &result)
	return result, nil
}

func main() {
	base := "http://120.77.251.18/api/v1"
	users := []struct{ name, pass, nick string }{
		{"admin", "123456", "管理员"},
		{"u1", "111111", "用户一"},
		{"u2", "111111", "用户二"},
	}

	fmt.Println("=== 注册 ===")
	for _, u := range users {
		resp, _ := request(base, "/register", map[string]any{
			"username": u.name,
			"password": u.pass,
			"nickname": u.nick,
		})
		fmt.Printf("%s: code=%.0f msg=%s\n", u.name, resp["code"].(float64), resp["message"])
	}

	fmt.Println("=== 登录 ===")
	for _, u := range users {
		resp, _ := request(base, "/login", map[string]any{
			"username": u.name,
			"password": u.pass,
		})
		fmt.Printf("%s: code=%.0f id=%.0f nick=%s\n", u.name, resp["code"].(float64), resp["user_id"].(float64), resp["nickname"])
	}
}
