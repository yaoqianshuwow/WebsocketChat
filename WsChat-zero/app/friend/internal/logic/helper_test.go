package logic

import (
	"testing"

	"github.com/your-org/ws-chat-zero/app/friend/internal/model"
)

func TestContactModelToProto(t *testing.T) {
	c := &model.Contact{
		Id:          1,
		UserId:      100,
		ContactId:   200,
		ContactType: 1,
		Status:      0,
	}

	pb := contactModelToProto(c)
	if pb == nil {
		t.Fatal("contactModelToProto returned nil")
	}
	if pb.Id != 1 {
		t.Errorf("expected Id=1, got %d", pb.Id)
	}
	if pb.UserId != 100 {
		t.Errorf("expected UserId=100, got %d", pb.UserId)
	}
	if pb.ContactId != 200 {
		t.Errorf("expected ContactId=200, got %d", pb.ContactId)
	}
}

func TestGroupModelToProto(t *testing.T) {
	g := &model.GroupInfo{
		Id:          1,
		Name:        "测试群",
		OwnerId:     100,
		MemberCount: 5,
		Status:      0,
	}

	pb := groupModelToProto(g)
	if pb == nil {
		t.Fatal("groupModelToProto returned nil")
	}
	if pb.Name != "测试群" {
		t.Errorf("expected Name=测试群, got %s", pb.Name)
	}
	if pb.MemberCount != 5 {
		t.Errorf("expected MemberCount=5, got %d", pb.MemberCount)
	}
}

func TestSessionModelToProto(t *testing.T) {
	s := &model.Session{
		Id:          1,
		UserId:      100,
		PeerId:      200,
		SessionType: 1,
		SessionName: "test",
	}

	pb := sessionModelToProto(s)
	if pb == nil {
		t.Fatal("sessionModelToProto returned nil")
	}
	if pb.SessionName != "test" {
		t.Errorf("expected SessionName=test, got %s", pb.SessionName)
	}
}
