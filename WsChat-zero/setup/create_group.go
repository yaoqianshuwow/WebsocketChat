package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func req(base, path string, token string, body map[string]any) map[string]any {
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
	base := "http://localhost:8888/api/v1"

	// 登录三个用户
	tokens := map[string]string{}
	for _, u := range []struct{ name, pass string }{{"u1", "111111"}, {"admin", "123456"}, {"u2", "111111"}} {
		r := req(base, "/login", "", map[string]any{"username": u.name, "password": u.pass})
		tokens[u.name] = r["token"].(string)
		fmt.Printf("%s logged in, id=%.0f\n", u.name, r["user_id"])
	}

	// u1 建群一家人
	r := req(base, "/group/createGroup", tokens["u1"], map[string]any{"groupName": "一家人"})
	fmt.Printf("create group: code=%.0f name=%s id=%.0f\n", r["code"], r["name"], r["group_id"])
	groupId := int(r["group_id"].(float64))

	// u1 加 admin 和 u2 入群（通过 joinGroup 不行，手动 INSERT）
	// 先让 u2 和 admin 加群
	for _, who := range []string{"admin", "u2"} {
		r := req(base, "/group/joinGroup", tokens[who], map[string]any{"groupId": float64(groupId)})
		fmt.Printf("%s join group: code=%.0f msg=%s\n", who, r["code"], r["message"])
	}

	// 查群成员
	r = req(base, "/group/getGroupMemberList", tokens["u1"], map[string]any{"groupId": float64(groupId)})
	members := r["memberList"].([]interface{})
	fmt.Printf("\n=== 一家人 群成员(%d人) ===\n", len(members))
	for _, m := range members {
		member := m.(map[string]any)
		fmt.Printf("  userId=%.0f nick=%s role=%.0f\n", member["userId"], member["nickname"], member["role"])
	}
}
