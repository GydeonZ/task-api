package handler

import (
	"strconv"

	"github.com/GydeonZ/task-api/pkg/apperrors"
)

// maxUint is the maximum value representable by uint on this platform.
const maxUint = ^uint(0)

// parseID safely parses a URL parameter string to uint,
// validating the value fits within the platform's uint size.
func parseID(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, apperrors.ErrBadRequest
	}
	if val > uint64(maxUint) {
		return 0, apperrors.ErrBadRequest
	}
	return uint(val), nil
}
