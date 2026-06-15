package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/segmentio/kafka-go"
	"github.com/your-org/ws-chat-zero/app/msg-store/internal/model"
	"github.com/your-org/ws-chat-zero/app/msg-store/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type MessagePayload struct {
	MsgId        int64  `json:"msg_id"`
	SenderId     int64  `json:"sender_id"`
	ReceiverId   int64  `json:"receiver_id"`
	ChatType     int32  `json:"chat_type"`
	MsgType      int32  `json:"msg_type"`
	Content      string `json:"content"`
	FileUrl      string `json:"file_url"`
	FileSize     int64  `json:"file_size"`
	FileName     string `json:"file_name"`
	SessionId    int64  `json:"session_id"`
	CreatedAt    int64  `json:"created_at"`
	SenderName   string `json:"sender_name"`
	SenderAvatar string `json:"sender_avatar"`
}

func StartMessageConsumer(ctx *svc.ServiceContext) {
	if err := ensureTopic(ctx.Config.Kafka.Brokers, ctx.Config.Kafka.ChatTopic); err != nil {
		logx.Errorf("ensure topic failed: %v", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     ctx.Config.Kafka.Brokers,
		Topic:       ctx.Config.Kafka.ChatTopic,
		GroupID:     ctx.Config.Kafka.GroupId,
		MinBytes:    10,
		MaxBytes:    10e6,
		MaxWait:     1 * time.Second,
		StartOffset: kafka.FirstOffset,
	})
	defer reader.Close()

	logx.Infof("Kafka consumer started, topic: %s, group: %s, brokers: %s",
		ctx.Config.Kafka.ChatTopic, ctx.Config.Kafka.GroupId, strings.Join(ctx.Config.Kafka.Brokers, ","))

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			logx.Errorf("kafka read error: %v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		var payload MessagePayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			logx.Errorf("unmarshal message error: %v", err)
			continue
		}

		if err := storeToMySQL(ctx, &payload); err != nil {
			logx.Errorf("store to MySQL error: %v", err)
		}

		if payload.MsgType == 1 {
			if err := storeToElasticsearch(ctx, &payload); err != nil {
				logx.Errorf("store to ES error: %v", err)
			}
		}

		logx.Infof("message stored: id=%d, type=%d, partition=%d, offset=%d", payload.MsgId, payload.MsgType, msg.Partition, msg.Offset)
	}
}

func ensureTopic(brokers []string, topic string) error {
	if len(brokers) == 0 {
		return fmt.Errorf("no kafka brokers configured")
	}

	var lastErr error
	for _, broker := range brokers {
		conn, err := kafka.Dial("tcp", broker)
		if err != nil {
			lastErr = err
			continue
		}

		controller, err := conn.Controller()
		if err != nil {
			lastErr = err
			_ = conn.Close()
			continue
		}

		controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, fmt.Sprintf("%d", controller.Port)))
		if err != nil {
			lastErr = err
			_ = conn.Close()
			continue
		}

		err = controllerConn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
		_ = controllerConn.Close()
		_ = conn.Close()
		if err != nil && !strings.Contains(strings.ToLower(err.Error()), "already exists") {
			lastErr = err
			continue
		}
		return nil
	}

	return lastErr
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
		SendName:   payload.SenderName,
		SendAvatar: payload.SenderAvatar,
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

	doc := map[string]any{
		"msg_id":      payload.MsgId,
		"sender_id":   payload.SenderId,
		"receiver_id": payload.ReceiverId,
		"chat_type":   payload.ChatType,
		"msg_type":    payload.MsgType,
		"content":     payload.Content,
		"session_id":  payload.SessionId,
		"created_at":  payload.CreatedAt,
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
