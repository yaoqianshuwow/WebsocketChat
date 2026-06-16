package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func req(base, path string, body map[string]any) (map[string]any, error) {
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

	fmt.Println("=== SERVER 注册 ===")
	for _, u := range users {
		r, e := req(base, "/register", map[string]any{
			"username": u.name, "password": u.pass, "nickname": u.nick,
		})
		if e != nil {
			fmt.Printf("%s: ERR %v\n", u.name, e)
			continue
		}
		fmt.Printf("%s: code=%.0f %v\n", u.name, r["code"].(float64), r["message"])
	}

	fmt.Println("=== SERVER 登录 ===")
	for _, u := range users {
		r, e := req(base, "/login", map[string]any{
			"username": u.name, "password": u.pass,
		})
		if e != nil {
			fmt.Printf("%s: ERR %v\n", u.name, e)
			continue
		}
		nick := ""
		if v, ok := r["nickname"]; ok && v != nil {
			nick = fmt.Sprintf("%v", v)
		}
		fmt.Printf("%s: ok=%.0f id=%.0f nick=%s\n", u.name, r["code"].(float64), r["user_id"].(float64), nick)
	}
}
