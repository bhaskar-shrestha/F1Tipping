package services

// ValidationError represents a client input validation error (HTTP 400)
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}
