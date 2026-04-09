package handler

import (
	"math/bits"
	"strconv"

	"github.com/GydeonZ/task-api/pkg/apperrors"
)

// parseID safely parses a URL parameter string to uint,
// validating the value fits within the platform's uint size.
func parseID(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, bits.UintSize)
	if err != nil {
		return 0, apperrors.ErrBadRequest
	}
	return uint(val), nil
}
