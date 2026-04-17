package errors


// AppError represents application-specific errors
type AppError struct {
	Code    string
	Message string
	Status  int
	Err     error
}

func (e *AppError) Error() string {
    return e.Message  // Hanya return message, ignore underlying error
}

// Common error codes and constructors
var (
	// Validation errors
	ErrInvalidInput = &AppError{
		Code:    "INVALID_INPUT",
		Message: "Invalid input provided",
		Status:  400,
	}

	// Not found errors
	ErrBoardNotFound = &AppError{
		Code:    "BOARD_NOT_FOUND",
		Message: "Board not found",
		Status:  404,
	}

	ErrListNotFound = &AppError{
		Code:    "LIST_NOT_FOUND",
		Message: "List not found",
		Status:  404,
	}

	ErrTaskNotFound = &AppError{
		Code:    "TASK_NOT_FOUND",
		Message: "Task not found",
		Status:  404,
	}

	// Conflict errors
	ErrBoardAlreadyExists = &AppError{
		Code:    "BOARD_ALREADY_EXISTS",
		Message: "Board already exists",
		Status:  409,
	}

	// Server errors
	ErrInternalServer = &AppError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
		Status:  500,
	}
)

// New creates a new AppError with a custom message
func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// NewWithErr creates a new AppError with an underlying error
func NewWithErr(code, message string, status int, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}
