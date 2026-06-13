package shared

import "fmt"

// DomainError is the base error type for all domain errors.
type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewDomainError(code, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

// Sentinel domain errors
var (
	ErrNotFound          = NewDomainError("NOT_FOUND", "resource not found")
	ErrAlreadyExists     = NewDomainError("ALREADY_EXISTS", "resource already exists")
	ErrInvalidInput      = NewDomainError("INVALID_INPUT", "invalid input data")
	ErrUnauthorized      = NewDomainError("UNAUTHORIZED", "not authorized")
	ErrForbidden         = NewDomainError("FORBIDDEN", "access forbidden")
	ErrTenantMismatch    = NewDomainError("TENANT_MISMATCH", "resource belongs to different tenant")
	ErrQuotaExceeded     = NewDomainError("QUOTA_EXCEEDED", "tenant quota exceeded")
	ErrInvalidStatus     = NewDomainError("INVALID_STATUS", "invalid status transition")
	ErrIPConflict        = NewDomainError("IP_CONFLICT", "IP address already allocated")
	ErrRackUnitConflict  = NewDomainError("RACK_UNIT_CONFLICT", "rack units already occupied")
	ErrAgentOffline      = NewDomainError("AGENT_OFFLINE", "agent is not connected")
	ErrLicenseExhausted  = NewDomainError("LICENSE_EXHAUSTED", "no available license seats")
)

// ValidationError carries field-level validation failures.
type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %v", e.Fields)
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Fields: map[string]string{field: message}}
}

func (e *ValidationError) Add(field, message string) *ValidationError {
	e.Fields[field] = message
	return e
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Fields) > 0
}
