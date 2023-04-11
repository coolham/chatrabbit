package httputil

import (
	"chatrabbit/config/ehandle"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/kataras/iris/v12/mvc"
)

func ErrCodeResp(err error) mvc.Response {
	var customErr = new(ehandle.CustomError)
	var result commons.ResultData
	if errors.As(err, &customErr) {
		result = commons.RespResult(int(customErr.Code), err.Error(), nil)

	} else {
		result = commons.RespResult(9999, err.Error(), nil)
	}
	content, _ := json.Marshal(result)
	return mvc.Response{
		ContentType: "application/json",
		Content:     content,
	}
}

func CodeResp(code int) mvc.Response {
	result := commons.RespResult(code, commons.ErrMsg(code), nil)

	content, _ := json.Marshal(result)
	return mvc.Response{
		ContentType: "application/json",
		Content:     content,
	}
}

func CodeRespWithData(code int, data interface{}) mvc.Response {
	result := commons.RespResult(code, commons.ErrMsg(code), data)
	content, _ := json.Marshal(result)
	return mvc.Response{
		ContentType: "application/json",
		Content:     content,
	}
}

func CodeRespWithMsgData(code int, msg string, data interface{}) mvc.Response {
	result := commons.RespResult(code, msg, data)
	content, _ := json.Marshal(result)
	return mvc.Response{
		ContentType: "application/json",
		Content:     content,
	}
}
