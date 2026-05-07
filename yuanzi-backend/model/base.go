package model

import (
	"crypto/rand"
	"fmt"
)

const (
	SUCCESS = 200

	ERROR                = 500
	ERROR_INVALID        = 400
	ERROR_NOT_AUTH       = 401
	ERROR_FORBID         = 403
	ERROR_NOT_FUND       = 404
	ERROR_QUOTA_EXCEEDED = 409
	ERROR_RATE_LIMITED   = 429

	ERROR_AUTH_CHECK_TOKEN_FAIL    = 10001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 10002
	ERROR_AUTH_TOKEN               = 10003
	ERROR_AUTH                     = 10004
)

var MsgFlags = map[int]string{
	SUCCESS:                        "ok",
	ERROR:                          "fail",
	ERROR_INVALID:                  "请求参数错误",
	ERROR_NOT_AUTH:                 "未认证",
	ERROR_FORBID:                   "无权限",
	ERROR_NOT_FUND:                 "资源不存在",
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "Token错误",
	ERROR_QUOTA_EXCEEDED:           "配额已用完",
	ERROR_RATE_LIMITED:             "请求过于频繁",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type ListResponse struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

// NewID 生成符合数据库 varchar(36) 主键格式的 UUID 风格字符串。
func NewID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16],
	)
}
