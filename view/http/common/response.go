package common

const (
	// ErrorCode 一般错误码
	ErrorCode = -1
	//成功返回码
	SuccessCode = 0
)

var (
	//SuccessCodeWithMessage 成功返回码与信息参数
	SuccessCodeWithMessage = codeWithMessage{Code: SuccessCode, Message: "Success"}

	//ParamErrorCodeWithMessage 参数错误返回码与返回信息
	ParamErrorCodeWithMessage = codeWithMessage{Code: -2, Message: "参数错误"}
	//MissingTokenCodeWithMessage 缺少用户token
	MissingTokenCodeWithMessage = codeWithMessage{Code: -3, Message: "当前用户未登录，请登录"}
	//OutOfPermissionRangeCodeWithMessage 请求参数超出权限范围错误返回码与返回信息
	OutOfPermissionRangeCodeWithMessage = codeWithMessage{Code: -4, Message: "请求参数超出权限范围或不存在，请调整"}
)

/*
*
返回码与返回信息
*/
type codeWithMessage struct {
	//返回码
	Code int
	//返回信息
	Message string
}

/*
返回实体类
*/
type Response struct {
	codeWithMessage
	//具体数据
	Data interface{}
}

/*
*

  - @description 自定义错误

  - @param

  - @return
    *
*/
func CustomErrorResponse(code int, message string) *Response {
	if code <= -1 {
		code = ErrorCode
	}
	return &Response{
		codeWithMessage: codeWithMessage{
			Code:    code,
			Message: message,
		},
		Data: nil,
	}
}

/*
*

  - @description 生成错误返回结构体

  - @param

  - @return
    *
*/
func ErrorResponse(message codeWithMessage) *Response {
	return &Response{message, nil}
}

/*
*

  - @description 成功respnse

  - @param

  - @return
    *
*/
func SuccessResponse(data interface{}) *Response {
	return &Response{SuccessCodeWithMessage, data}
}
