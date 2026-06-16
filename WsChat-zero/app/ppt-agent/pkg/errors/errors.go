package errors

import "fmt"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *AppError) WithMessage(msg string) *AppError {
	return &AppError{Code: e.Code, Message: msg}
}

var (
	SystemError    = &AppError{Code: 500, Message: "系统错误"}
	ParamsError    = &AppError{Code: 400, Message: "参数错误"}
	NotFoundError  = &AppError{Code: 404, Message: "未找到"}
	AuthError      = &AppError{Code: 401, Message: "未授权"}
)
