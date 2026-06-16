package vo

// PptGenerateRequest 生成PPT的请求
type PptGenerateRequest struct {
	Topic   string `json:"topic"`             // PPT主题
	Style   string `json:"style"`             // 风格（简约/商务/科技/创意等）
	Outline string `json:"outline,omitempty"` // 可选：自定义大纲
}

// PptGenerateResponse PPT生成响应
type PptGenerateResponse struct {
	TaskID  string `json:"taskId"`  // 任务ID
	Status  string `json:"status"`  // 状态
	Message string `json:"message"` // 消息
	Topic   string `json:"topic,omitempty"`
	Style   string `json:"style,omitempty"`
}

// PptDownloadResponse PPT下载响应
type PptDownloadResponse struct {
	FileURL  string `json:"fileUrl"`  // 下载URL
	FileName string `json:"fileName"` // 文件名
}

// SlideInfo 单页幻灯片信息
type SlideInfo struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	ImageURL string `json:"imageUrl,omitempty"`
	Chart    string `json:"chart,omitempty"`
}

// PptPreviewResponse PPT预览响应
type PptPreviewResponse struct {
	TaskID    string      `json:"taskId"`
	Status    string      `json:"status"`
	TotalPage int         `json:"totalPage"`
	Slides    []*SlideInfo `json:"slides"`
}

// PptTaskItem 任务列表项
type PptTaskItem struct {
	TaskID    string `json:"taskId"`
	Topic     string `json:"topic"`
	Style     string `json:"style"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}
