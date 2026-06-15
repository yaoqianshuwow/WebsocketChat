package logic

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/your-org/ws-chat-zero/app/user/internal/svc"
	"github.com/your-org/ws-chat-zero/app/user/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUsersLogic {
	return &SearchUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUsersLogic) SearchUsers(in *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	if in.Page <= 0 {
		in.Page = 1
	}
	if in.Size <= 0 {
		in.Size = 20
	}

	// 优先走 ES 搜索（含中文分词）
	if l.svcCtx.ES != nil {
		return l.searchViaES(in)
	}

	// 兜底：MySQL LIKE 搜索
	return l.searchViaMySQL(in)
}

func (l *SearchUsersLogic) searchViaES(in *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	from := (in.Page - 1) * in.Size
	index := l.svcCtx.Config.ES.Index
	if index == "" {
		index = "users"
	}

	query := map[string]any{
		"from": from,
		"size": in.Size,
		"query": map[string]any{
			"bool": map[string]any{
				"should": []map[string]any{
					{
						"multi_match": map[string]any{
							"query":  in.Keyword,
							"fields": []string{"username", "nickname"},
							"type":   "best_fields",
						},
					},
					{
						"multi_match": map[string]any{
							"query":  in.Keyword,
							"fields": []string{"username.autocomplete", "nickname.autocomplete"},
							"type":   "bool_prefix",
						},
					},
				},
			},
		},
		"sort": []map[string]any{{"_score": "desc"}},
	}

	body, _ := json.Marshal(query)
	resp, err := l.svcCtx.ES.Search(
		l.svcCtx.ES.Search.WithIndex(index),
		l.svcCtx.ES.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		logx.Errorf("ES search error: %v, fallback to MySQL", err)
		return l.searchViaMySQL(in)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		logx.Errorf("ES search error response: %s, fallback to MySQL", resp.String())
		return l.searchViaMySQL(in)
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source map[string]any `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logx.Errorf("ES decode error: %v, fallback to MySQL", err)
		return l.searchViaMySQL(in)
	}

	var data []*pb.UserInfo
	for _, hit := range result.Hits.Hits {
		user := hit.Source
		userId, _ := user["user_id"].(float64)
		var status int32
		if s, ok := user["status"].(float64); ok {
			status = int32(s)
		}
		nickname, _ := user["nickname"].(string)
		username, _ := user["username"].(string)
		avatar, _ := user["avatar"].(string)
		phone, _ := user["phone"].(string)

		data = append(data, &pb.UserInfo{
			Id:       int64(userId),
			Username: username,
			Nickname: nickname,
			Avatar:   avatar,
			Phone:    phone,
			Status:   status,
		})
	}

	return &pb.SearchUsersResponse{
		Code:    0,
		Message: "ok",
		Data:    data,
		Total:   result.Hits.Total.Value,
	}, nil
}

func (l *SearchUsersLogic) searchViaMySQL(in *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	var users []*pb.UserInfo
	var total int64

	query := l.svcCtx.DB.Model(&struct{}{}).Table("user_info").
		Where("username LIKE ? OR nickname LIKE ?", "%"+in.Keyword+"%", "%"+in.Keyword+"%")
	query.Count(&total)
	query.Offset(int((in.Page - 1) * in.Size)).Limit(int(in.Size)).Find(&users)

	for _, u := range users {
		if u.Nickname == "" {
			u.Nickname = u.Username
		}
	}

	return &pb.SearchUsersResponse{Code: 0, Message: "ok", Data: users, Total: total}, nil
}
