package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is the standard JSON error envelope returned by all API errors.
type ErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// ValidationErrorResponse carries per-field validation failures.
type ValidationErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields"`
}

func BadRequest(c *gin.Context, code, msg string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{Code: code, Message: msg})
}

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, ErrorResponse{Code: "UNAUTHORIZED", Message: msg})
}

func Forbidden(c *gin.Context, msg string) {
	c.JSON(http.StatusForbidden, ErrorResponse{Code: "FORBIDDEN", Message: msg})
}

func NotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, ErrorResponse{
		Code:    "NOT_FOUND",
		Message: resource + " not found",
	})
}

func Conflict(c *gin.Context, msg string) {
	c.JSON(http.StatusConflict, ErrorResponse{Code: "CONFLICT", Message: msg})
}

func UnprocessableEntity(c *gin.Context, fields map[string]string) {
	c.JSON(http.StatusUnprocessableEntity, ValidationErrorResponse{
		Code:    "VALIDATION_ERROR",
		Message: "validation failed",
		Fields:  fields,
	})
}

func TooManyRequests(c *gin.Context) {
	c.JSON(http.StatusTooManyRequests, ErrorResponse{
		Code:    "RATE_LIMITED",
		Message: "too many requests, slow down",
	})
}

func Internal(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Code:    "INTERNAL_ERROR",
		Message: "an internal error occurred",
	})
}

// AbortWithError writes the error and aborts the Gin chain.
func AbortUnauthorized(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Code: "UNAUTHORIZED", Message: msg})
}

func AbortForbidden(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{Code: "FORBIDDEN", Message: msg})
}

// StatusCode maps domain error codes to HTTP status codes.
func StatusCode(domainCode string) int {
	switch domainCode {
	case "NOT_FOUND":
		return http.StatusNotFound
	case "ALREADY_EXISTS", "IP_CONFLICT", "RACK_UNIT_CONFLICT":
		return http.StatusConflict
	case "INVALID_INPUT", "VALIDATION_ERROR":
		return http.StatusBadRequest
	case "UNAUTHORIZED":
		return http.StatusUnauthorized
	case "FORBIDDEN":
		return http.StatusForbidden
	case "QUOTA_EXCEEDED", "LICENSE_EXHAUSTED":
		return http.StatusPaymentRequired
	default:
		return http.StatusInternalServerError
	}
}
