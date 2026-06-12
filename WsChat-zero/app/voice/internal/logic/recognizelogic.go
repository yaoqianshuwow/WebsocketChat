package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/your-org/ws-chat-zero/app/voice/internal/svc"
	"github.com/your-org/ws-chat-zero/app/voice/pb/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecognizeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecognizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecognizeLogic {
	return &RecognizeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BaiduTokenManager 管理百度 access token
type BaiduTokenManager struct {
	mu        sync.RWMutex
	apiKey    string
	secretKey string
	token     string
	expireAt  time.Time
}

func NewBaiduTokenManager(apiKey, secretKey string) *BaiduTokenManager {
	return &BaiduTokenManager{
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

func (m *BaiduTokenManager) GetToken() (string, error) {
	m.mu.RLock()
	if m.token != "" && time.Now().Before(m.expireAt) {
		m.mu.RUnlock()
		return m.token, nil
	}
	m.mu.RUnlock()
	return m.refresh()
}

func (m *BaiduTokenManager) refresh() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查
	if m.token != "" && time.Now().Before(m.expireAt) {
		return m.token, nil
	}

	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s",
		m.apiKey, m.secretKey)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", fmt.Errorf("获取百度token失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("百度token获取失败: %s", result.Error)
	}

	m.token = result.AccessToken
	m.expireAt = time.Now().Add(time.Duration(result.ExpiresIn-60) * time.Second)
	return m.token, nil
}

// BaiduAsrResponse 百度语音识别响应
type BaiduAsrResponse struct {
	ErrNo int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Result []string `json:"result"`
}

func (l *RecognizeLogic) Recognize(in *pb.RecognizeRequest) (*pb.RecognizeResponse, error) {
	// 获取音频格式
	format := in.AudioFormat
	if format == "" {
		format = "pcm"
	}
	sampleRate := in.SampleRate
	if sampleRate <= 0 {
		sampleRate = 16000
	}

	tokenManager := NewBaiduTokenManager(
		l.svcCtx.Config.BaiduAsr.ApiKey,
		l.svcCtx.Config.BaiduAsr.SecretKey,
	)

	token, err := tokenManager.GetToken()
	if err != nil {
		logx.Errorf("获取百度token失败: %v", err)
		// 返回模拟识别结果（当百度API不可用时）
		return l.fallbackRecognize(in)
	}

	// 百度语音识别 API
	// https://ai.baidu.com/ai-doc/SPEECH/Vk38lxily
	url := fmt.Sprintf("https://vop.baidu.com/server_api?access_token=%s", token)

	speechLen := len(in.AudioData)
	payload := map[string]interface{}{
		"format":  format,
		"rate":    sampleRate,
		"channel": 1,
		"cuid":    fmt.Sprintf("user_%d", in.UserId),
		"token":   token,
		"len":     speechLen,
		"speech": map[string]interface{}{
			"audio": in.AudioData,
		},
	}

	bodyBytes, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		logx.Errorf("百度ASR请求失败: %v", err)
		return l.fallbackRecognize(in)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var asrResp BaiduAsrResponse
	if err := json.Unmarshal(respBody, &asrResp); err != nil {
		return l.fallbackRecognize(in)
	}

	if asrResp.ErrNo != 0 {
		logx.Errorf("百度ASR识别失败: err_no=%d, err_msg=%s", asrResp.ErrNo, asrResp.ErrMsg)
		return &pb.RecognizeResponse{
			Code:    1,
			Message: fmt.Sprintf("识别失败: %s", asrResp.ErrMsg),
			UserId:  in.UserId,
			MsgId:   in.MsgId,
		}, nil
	}

	text := ""
	if len(asrResp.Result) > 0 {
		text = asrResp.Result[0]
	}

	return &pb.RecognizeResponse{
		Code:    0,
		Message: "ok",
		Text:    text,
		UserId:  in.UserId,
		MsgId:   in.MsgId,
	}, nil
}

// fallbackRecognize 当百度API不可用时的回退方案（模拟识别）
func (l *RecognizeLogic) fallbackRecognize(in *pb.RecognizeRequest) (*pb.RecognizeResponse, error) {
	logx.Infof("使用模拟语音识别: user_id=%d, audio_size=%d, format=%s",
		in.UserId, len(in.AudioData), in.AudioFormat)

	// TODO: 当百度API配置完善后，此回退应移除
	return &pb.RecognizeResponse{
		Code:    0,
		Message: "ok",
		Text:    "[语音识别暂不可用，请配置百度API]",
		UserId:  in.UserId,
		MsgId:   in.MsgId,
	}, nil
}
