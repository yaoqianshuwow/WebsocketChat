package myutils

import (
	"fmt"
	"time"
)

// GenerateShortID 生成短ID
func GenerateShortID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
