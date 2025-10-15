package common

import (
	"errors"
	"fmt"
)

// CommonError represents a structured error type for the common package
type CommonError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

// ErrorCode represents different types of errors
type ErrorCode string

const (
	// Directory operation errors
	ErrDirCreation ErrorCode = "DIRECTORY_CREATION_FAILED"
	ErrDirAccess   ErrorCode = "DIRECTORY_ACCESS_FAILED"

	// File operation errors
	ErrFileNotFound ErrorCode = "FILE_NOT_FOUND"
	ErrFileRead     ErrorCode = "FILE_READ_FAILED"
	ErrFileWrite    ErrorCode = "FILE_WRITE_FAILED"
	ErrFileCreate   ErrorCode = "FILE_CREATE_FAILED"

	// Network operation errors
	ErrNetworkRequest  ErrorCode = "NETWORK_REQUEST_FAILED"
	ErrNetworkResponse ErrorCode = "NETWORK_RESPONSE_FAILED"

	// Template and configuration errors
	ErrTemplateLoad  ErrorCode = "TEMPLATE_LOAD_FAILED"
	ErrTemplateParse ErrorCode = "TEMPLATE_PARSE_FAILED"
	ErrConfigInvalid ErrorCode = "CONFIG_INVALID"

	// Subscription errors
	ErrSubscriptionLoad  ErrorCode = "SUBSCRIPTION_LOAD_FAILED"
	ErrSubscriptionParse ErrorCode = "SUBSCRIPTION_PARSE_FAILED"

	// Regex errors
	ErrRegexCompile ErrorCode = "REGEX_COMPILE_FAILED"
	ErrRegexInvalid ErrorCode = "REGEX_INVALID"

	// Database errors
	ErrDatabaseConnect ErrorCode = "DATABASE_CONNECTION_FAILED"
	ErrDatabaseQuery   ErrorCode = "DATABASE_QUERY_FAILED"
	ErrRecordNotFound  ErrorCode = "RECORD_NOT_FOUND"

	// Validation errors
	ErrValidation   ErrorCode = "VALIDATION_FAILED"
	ErrInvalidInput ErrorCode = "INVALID_INPUT"
)

// Error returns the string representation of the error
func (e *CommonError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *CommonError) Unwrap() error {
	return e.Cause
}

// Is allows error comparison
func (e *CommonError) Is(target error) bool {
	if t, ok := target.(*CommonError); ok {
		return e.Code == t.Code
	}
	return false
}

// NewError creates a new CommonError
func NewError(code ErrorCode, message string, cause error) *CommonError {
	return &CommonError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewSimpleError creates a new CommonError without a cause
func NewSimpleError(code ErrorCode, message string) *CommonError {
	return &CommonError{
		Code:    code,
		Message: message,
	}
}

// Convenience constructors for common error types

// Directory errors
func NewDirCreationError(dirPath string, cause error) *CommonError {
	return NewError(ErrDirCreation, fmt.Sprintf("failed to create directory: %s", dirPath), cause)
}

func NewDirAccessError(dirPath string, cause error) *CommonError {
	return NewError(ErrDirAccess, fmt.Sprintf("failed to access directory: %s", dirPath), cause)
}

// File errors
func NewFileNotFoundError(filePath string) *CommonError {
	return NewSimpleError(ErrFileNotFound, fmt.Sprintf("file not found: %s", filePath))
}

func NewFileReadError(filePath string, cause error) *CommonError {
	return NewError(ErrFileRead, fmt.Sprintf("failed to read file: %s", filePath), cause)
}

func NewFileWriteError(filePath string, cause error) *CommonError {
	return NewError(ErrFileWrite, fmt.Sprintf("failed to write file: %s", filePath), cause)
}

func NewFileCreateError(filePath string, cause error) *CommonError {
	return NewError(ErrFileCreate, fmt.Sprintf("failed to create file: %s", filePath), cause)
}

// Network errors
func NewNetworkRequestError(url string, cause error) *CommonError {
	return NewError(ErrNetworkRequest, fmt.Sprintf("network request failed for URL: %s", url), cause)
}

func NewNetworkResponseError(message string, cause error) *CommonError {
	return NewError(ErrNetworkResponse, message, cause)
}

// Template errors
func NewTemplateLoadError(template string, cause error) *CommonError {
	return NewError(ErrTemplateLoad, fmt.Sprintf("failed to load template: %s", template), cause)
}

func NewTemplateParseError(data []byte, cause error) *CommonError {
	return NewError(ErrTemplateParse, fmt.Sprintf("failed to parse template: %s", data), cause)
}

// Subscription errors
func NewSubscriptionLoadError(url string, cause error) *CommonError {
	return NewError(ErrSubscriptionLoad, fmt.Sprintf("failed to load subscription: %s", url), cause)
}

func NewSubscriptionParseError(data []byte, cause error) *CommonError {
	return NewError(ErrSubscriptionParse, fmt.Sprintf("failed to parse subscription: %s", string(data)), cause)
}

// Regex errors
func NewRegexCompileError(pattern string, cause error) *CommonError {
	return NewError(ErrRegexCompile, fmt.Sprintf("failed to compile regex pattern: %s", pattern), cause)
}

func NewRegexInvalidError(paramName string, cause error) *CommonError {
	return NewError(ErrRegexInvalid, fmt.Sprintf("invalid regex in parameter: %s", paramName), cause)
}

// Database errors
func NewDatabaseConnectError(cause error) *CommonError {
	return NewError(ErrDatabaseConnect, "failed to connect to database", cause)
}

func NewRecordNotFoundError(recordType string, id string) *CommonError {
	return NewSimpleError(ErrRecordNotFound, fmt.Sprintf("%s not found: %s", recordType, id))
}

// Validation errors
func NewValidationError(field string, message string) *CommonError {
	return NewSimpleError(ErrValidation, fmt.Sprintf("validation failed for %s: %s", field, message))
}

func NewInvalidInputError(paramName string, value string) *CommonError {
	return NewSimpleError(ErrInvalidInput, fmt.Sprintf("invalid input for parameter %s: %s", paramName, value))
}

// IsErrorCode checks if an error has a specific error code
func IsErrorCode(err error, code ErrorCode) bool {
	var commonErr *CommonError
	if errors.As(err, &commonErr) {
		return commonErr.Code == code
	}
	return false
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) (ErrorCode, bool) {
	var commonErr *CommonError
	if errors.As(err, &commonErr) {
		return commonErr.Code, true
	}
	return "", false
}
