package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/your-org/ws-chat-zero/app/msg-store/internal/model"
	"github.com/your-org/ws-chat-zero/app/msg-store/internal/svc"
)

// MessagePayload Kafka 消息体
type MessagePayload struct {
	MsgId      int64  `json:"msg_id"`
	SenderId   int64  `json:"sender_id"`
	ReceiverId int64  `json:"receiver_id"`
	ChatType   int32  `json:"chat_type"`
	MsgType    int32  `json:"msg_type"`
	Content    string `json:"content"`
	FileUrl    string `json:"file_url"`
	FileSize   int64  `json:"file_size"`
	FileName   string `json:"file_name"`
	SessionId  int64  `json:"session_id"`
	CreatedAt  int64  `json:"created_at"`
}

func StartMessageConsumer(ctx *svc.ServiceContext) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   ctx.Config.Kafka.Brokers,
		Topic:     ctx.Config.Kafka.ChatTopic,
		GroupID:   ctx.Config.Kafka.GroupId,
		MinBytes:  10,
		MaxBytes:  10e6, // 10MB
		MaxWait:   1 * time.Second,
	})

	logx.Infof("Kafka consumer started, topic: %s, group: %s",
		ctx.Config.Kafka.ChatTopic, ctx.Config.Kafka.GroupId)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			logx.Errorf("kafka read error: %v", err)
			continue
		}

		var payload MessagePayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			logx.Errorf("unmarshal message error: %v", err)
			continue
		}

		// 始终存储到 MySQL
		if err := storeToMySQL(ctx, &payload); err != nil {
			logx.Errorf("store to MySQL error: %v", err)
		}

		// 文本消息额外存储到 ES
		if payload.MsgType == 1 {
			if err := storeToElasticsearch(ctx, &payload); err != nil {
				logx.Errorf("store to ES error: %v", err)
			}
		}

		logx.Infof("message stored: id=%d, type=%d", payload.MsgId, payload.MsgType)
	}
}

func storeToMySQL(ctx *svc.ServiceContext, payload *MessagePayload) error {
	msg := model.Message{
		SenderId:   payload.SenderId,
		ReceiverId: payload.ReceiverId,
		ChatType:   payload.ChatType,
		MsgType:    payload.MsgType,
		Content:    payload.Content,
		FileUrl:    payload.FileUrl,
		FileSize:   payload.FileSize,
		FileName:   payload.FileName,
		Status:     0,
		SessionId:  payload.SessionId,
		CreatedAt:  time.Unix(payload.CreatedAt, 0),
	}
	if err := ctx.DB.Create(&msg).Error; err != nil {
		return fmt.Errorf("MySQL insert failed: %w", err)
	}
	logx.Infof("message saved to MySQL: msg_id=%d, table_id=%d", payload.MsgId, msg.Id)
	return nil
}

func storeToElasticsearch(ctx *svc.ServiceContext, payload *MessagePayload) error {
	if len(ctx.Config.ES.Addresses) == 0 {
		logx.Infof("ES addresses not configured, skip ES storage")
		return nil
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: ctx.Config.ES.Addresses,
	})
	if err != nil {
		return fmt.Errorf("ES client init failed: %w", err)
	}

	doc := map[string]interface{}{
		"msg_id":     payload.MsgId,
		"sender_id":  payload.SenderId,
		"receiver_id": payload.ReceiverId,
		"chat_type":  payload.ChatType,
		"msg_type":   payload.MsgType,
		"content":    payload.Content,
		"session_id": payload.SessionId,
		"created_at": payload.CreatedAt,
	}

	body, _ := json.Marshal(doc)
	resp, err := es.Index(
		ctx.Config.ES.Index,
		bytes.NewReader(body),
		es.Index.WithDocumentID(fmt.Sprintf("%d", payload.MsgId)),
	)
	if err != nil {
		return fmt.Errorf("ES index failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return fmt.Errorf("ES index error: %s", resp.String())
	}

	logx.Infof("message indexed to ES: msg_id=%d", payload.MsgId)
	return nil
}
