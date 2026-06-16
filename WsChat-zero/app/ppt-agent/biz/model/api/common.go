package api

import (
	pkg "ppt-agent/pkg/errors"
)

// Response 通用响应结构
type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// NewSuccessResponse 成功响应
func NewSuccessResponse[T any](data T) *Response[T] {
	return &Response[T]{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// NewErrorResponse 错误响应
func NewErrorResponse[T any](err *pkg.AppError) *Response[T] {
	return &Response[T]{
		Code:    err.Code,
		Message: err.Error(),
	}
}

// PageRequest 分页请求
type PageRequest struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

// PageResponse 分页响应
type PageResponse[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}
