package http

import "encoding/json"

type BaseResp struct {
	Code   int         `json:"code"`   //状态码
	Status int         `json:"status"` //兼容服务端
	Msg    string      `json:"msg"`    //错误码
	Data   interface{} `json:"data"`   //结果
}

const http_ok = "ok"

func NewBaseResp(data interface{}) *BaseResp {
	return &BaseResp{
		Code:   0,
		Status: 0,
		Msg:    http_ok,
		Data:   data,
	}
}

func NewFailResp(errMsg string) *BaseResp {
	return &BaseResp{
		Code:   -1,
		Status: -1,
		Msg:    errMsg,
	}
}

func NewFailRespBytes(errMsg string) []byte {
	resp := &BaseResp{
		Code:   -1,
		Status: -1,
		Msg:    errMsg,
	}
	bb, _ := json.Marshal(resp)
	return bb
}
