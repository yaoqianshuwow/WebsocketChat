package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func req(base, path, token string, body map[string]any) map[string]any {
	data, _ := json.Marshal(body)
	r, _ := http.NewRequest("POST", base+path, bytes.NewReader(data))
	r.Header.Set("Content-Type", "application/json")
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	resp, _ := http.DefaultClient.Do(r)
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	var result map[string]any
	json.Unmarshal(raw, &result)
	return result
}

func main() {
	base := os.Getenv("SETUP_URL")
	if base == "" {
		base = "http://localhost:8888/api/v1"
	}
	fmt.Printf("Target: %s\n\n", base)

	// Step 1: Register 3 accounts
	fmt.Println("=== 注册 ===")
	users := []struct{ name, pass, nick string }{
		{"admin", "123456", "管理员"},
		{"u1", "111111", "用户一"},
		{"u2", "111111", "用户二"},
	}
	for _, u := range users {
		r := req(base, "/register", "", map[string]any{
			"username": u.name, "password": u.pass, "nickname": u.nick,
		})
		fmt.Printf("%s (%s): code=%.0f %s\n", u.name, u.nick, r["code"], r["message"])
	}

	// Step 2: Login all
	fmt.Println("\n=== 登录 ===")
	tokens := map[string]int64{}
	for _, u := range users {
		r := req(base, "/login", "", map[string]any{"username": u.name, "password": u.pass})
		tokens[u.name] = int64(r["user_id"].(float64))
		fmt.Printf("%s: id=%.0f\n", u.name, r["user_id"])
	}

	// Step 3: u1 creates group
	fmt.Println("\n=== 建群 ===")
	u1Token := getToken(users, tokens)
	r := req(base, "/group/createGroup", u1Token, map[string]any{"groupName": "一家人"})
	gid := int64(r["group_id"].(float64))
	fmt.Printf("一家人: id=%d code=%.0f\n", gid, r["code"])

	// Step 4: admin and u2 join group
	for _, name := range []string{"admin", "u2"} {
		t := getTokenByName(name, users, tokens)
		r := req(base, "/group/joinGroup", t, map[string]any{"groupId": float64(gid)})
		fmt.Printf("%s join: code=%.0f %s\n", name, r["code"], r["message"])
	}

	// Step 5: Show group members
	r = req(base, "/group/getGroupMemberList", u1Token, map[string]any{"groupId": float64(gid)})
	members := r["memberList"].([]interface{})
	fmt.Printf("\n=== 一家人 成员(%d人) ===\n", len(members))
	for _, m := range members {
		member := m.(map[string]any)
		fmt.Printf("  userId=%.0f nick=%s role=%.0f\n", member["userId"], member["nickname"], member["role"])
	}

	fmt.Println("\n=== 完成 ===")
}

func getToken(users []struct{ name, pass, nick string }, tokens map[string]int64) string {
	for _, u := range users {
		if u.name == "u1" {
			for name, token := range tokens {
				if name == "u1" {
					return req("http://localhost:8888/api/v1", "/login", "", map[string]any{"username": "u1", "password": "111111"})["token"].(string)
				}
			}
		}
	}
	return ""
}

func getTokenByName(name string, users []struct{ name, pass, nick string }, tokens map[string]int64) string {
	for _, u := range users {
		if u.name == name {
			r := req("http://localhost:8888/api/v1", "/login", "", map[string]any{"username": name, "password": u.pass})
			return r["token"].(string)
		}
	}
	return ""
}
