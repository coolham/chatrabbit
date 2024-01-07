package response

import (
	"chatrabbit/apperrors"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/kataras/iris/v12/mvc"
)

// APIResponse is a common structure for API responses.
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewAPIResponse creates a new API response with the given code, message, and data.
func NewAPIResponse(code int, message string, data interface{}) APIResponse {
	return APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Preset responses for common scenarios
var (
	ResponseOK            = NewAPIResponse(0, "OK", nil)
	ResponseBadRequest    = NewAPIResponse(400, "Bad Request", nil)
	ResponseUnauthorized  = NewAPIResponse(401, "Unauthorized", nil)
	ResponseInternalError = NewAPIResponse(500, "Internal Server Error", nil)
)

// ErrCodeResp processes the given error and returns a mvc.Response with HTTP status 200
// and the appropriate error code and message in the JSON body.
func ErrCodeResp(err error) mvc.Response {
	var customErr = new(apperrors.CustomError)
	var result APIResponse
	if errors.As(err, &customErr) {
		msg := customErr.Message + ", " + customErr.Details
		result = NewAPIResponse(int(customErr.Code), msg, nil)
	} else {
		result = NewAPIResponse(9999, err.Error(), nil)
	}

	content, _ := json.Marshal(result)

	return mvc.Response{
		ContentType: "application/json",
		Content:     content,
		Code:        200, // Always return HTTP status 200
	}
}

// DataResp returns a mvc.Response with HTTP status 200, code 0, and the provided data in the JSON body.
func DataResp(data interface{}) mvc.Response {
	result := NewAPIResponse(0, "OK", data)
	content, _ := json.Marshal(result)

	return mvc.Response{
		ContentType: "application/json",
		Content:     content,
		Code:        200,
	}
}
