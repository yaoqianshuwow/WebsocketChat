package logic

import (
	"testing"

	"github.com/your-org/ws-chat-zero/app/msg-forward/internal/model"
)

func TestMessageModelToProto(t *testing.T) {
	m := &model.Message{
		Id:         1,
		SenderId:   100,
		ReceiverId: 200,
		ChatType:   1,
		MsgType:    1,
		Content:    "hello",
		SessionId:  10,
	}

	pb := messageModelToProto(m)
	if pb == nil {
		t.Fatal("messageModelToProto returned nil")
	}
	if pb.Content != "hello" {
		t.Errorf("expected Content=hello, got %s", pb.Content)
	}
	if pb.SenderId != 100 {
		t.Errorf("expected SenderId=100, got %d", pb.SenderId)
	}
}

func TestMessageTableName(t *testing.T) {
	m := &model.Message{}
	if m.TableName() != "message" {
		t.Errorf("expected table=message, got %s", m.TableName())
	}
}
