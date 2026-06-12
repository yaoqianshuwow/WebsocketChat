package logic

import (
	"testing"

	"github.com/your-org/ws-chat-zero/app/user/internal/model"
)

func TestModelToProto(t *testing.T) {
	u := &model.UserInfo{
		Id:       1,
		Username: "testuser",
		Nickname: "测试用户",
		Status:   0,
		Role:     0,
	}

	pb := modelToProto(u)
	if pb == nil {
		t.Fatal("modelToProto returned nil")
	}
	if pb.Id != 1 {
		t.Errorf("expected Id=1, got %d", pb.Id)
	}
	if pb.Username != "testuser" {
		t.Errorf("expected Username=testuser, got %s", pb.Username)
	}
	if pb.Nickname != "测试用户" {
		t.Errorf("expected Nickname=测试用户, got %s", pb.Nickname)
	}
	if pb.Status != 0 {
		t.Errorf("expected Status=0, got %d", pb.Status)
	}
}

func TestModelToProtoEmpty(t *testing.T) {
	u := &model.UserInfo{}
	pb := modelToProto(u)
	if pb == nil {
		t.Fatal("modelToProto returned nil for empty model")
	}
}
