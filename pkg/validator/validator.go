package validator

import (
	"regexp"
	"strings"
)

// ValidateEmail checks if the email is valid
func ValidateEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ValidatePassword checks if the password meets the criteria
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}

	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return false, "Password must contain at least one number"
	}

	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		return false, "Password must contain at least one special character"
	}

	return true, ""
}

// ValidateUsername checks if the username is valid
func ValidateUsername(username string) (bool, string) {
	username = strings.TrimSpace(username)

	if len(username) < 5 {
		return false, "Username must be at least 5 characters long"
	}
	if len(username) > 15 {
		return false, "Username must be no more than 15 characters long"
	}

	// Only alphanumeric characters and underscores allowed
	validChars := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
    if !validChars.MatchString(username) {
        return false, "Username can only contain letters, numbers, and underscores"
    }

    return true, ""
}