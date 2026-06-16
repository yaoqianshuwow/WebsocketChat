package snowflake

import (
	"github.com/sony/sonyflake"
	"strconv"
)

var sf *sonyflake.Sonyflake

func init() {
	sf = sonyflake.NewSonyflake(sonyflake.Settings{})
}

// NextID 生成唯一ID
func NextID() (uint64, error) {
	if sf == nil {
		sf = sonyflake.NewSonyflake(sonyflake.Settings{})
	}
	return sf.NextID()
}

// NextStrID 生成字符串格式的唯一ID
func NextStrID() string {
	id, err := NextID()
	if err != nil {
		return ""
	}
	return strconv.FormatUint(id, 10)
}
