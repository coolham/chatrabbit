package apperrors

import "github.com/pkg/errors"

//go:generate stringer -type ErrCode -linecomment

// ErrCode type for error codes
type ErrCode int

const (
	// sample error codes
	OK ErrCode = 0 // OK

	// User service error codes
	UserNotFound ErrCode = 1001 // 用户不存在

)

// CustomError represents a structured error for the API.
type CustomError struct {
	Code    ErrCode
	Message string
	// Optionally, add more fields for detailed error info or context
	Details string
}

// Error method to comply with the error interface.
func (e CustomError) Error() string {
	// return fmt.Sprintf("%s, %s", e.Message, e.Details)
	return e.Code.String()
}

// NewCustomError creates a new CustomError based on the error code.
// Additional details can be provided for richer error information.
func NewCustomError(code ErrCode, details string) error {
	// This uses the auto-generated string mapping for error messages
	// return CustomError{
	// 	Code:    code,
	// 	Message: code.String(),
	// 	Details: details,
	// }
	// 初次调用得用Wrap方法，进行实例化
	return errors.Wrap(&CustomError{
		Code:    code,
		Message: code.String(),
		Details: details,
	}, details)
}

// // An example of using the new error structure
// // This function is just for illustration and will not actually run.
// func ExampleAPIFunction() error {
// 	// ... some code
// 	return NewCustomError(BadRequest, "The provided email format is invalid.")
// }
