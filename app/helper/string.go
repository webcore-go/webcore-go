package helper

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// FormatDuration formats a duration as a human-readable string
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}

// GetClientIP gets the client IP address from a request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to remote address
	return r.RemoteAddr
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(uuid string) bool {
	if len(uuid) != 36 {
		return false
	}

	// Check UUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	pattern := `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`
	matched, _ := regexp.MatchString(pattern, uuid)
	return matched
}

// Truncate truncates a string to a specified length
func Truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

// Capitalize capitalizes the first letter of a string
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// ToSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	s = regexp.MustCompile(`[A-Z][a-z0-9]*`).ReplaceAllStringFunc(s, func(match string) string {
		return "_" + strings.ToLower(match)
	})
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return strings.TrimLeft(s, "_")
}

// ToCamelCase converts a string to camelCase
func ToCamelCase(s string) string {
	words := strings.Split(s, "_")
	if len(words) == 0 {
		return s
	}

	result := words[0]
	for i := 1; i < len(words); i++ {
		result += Capitalize(words[i])
	}

	return result
}

// GenerateID generates a unique ID
func GenerateID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GenerateUUID generates a UUID
func GenerateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

// HashPassword hashes a password (placeholder - implement actual hashing)
func HashPassword(password string) string {
	// In a real implementation, use bcrypt or similar
	return password
}

// ValidateEmail validates an email address
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePhone validates a phone number
func ValidatePhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^[\+]?[1-9][\d]{0,15}$`)
	return phoneRegex.MatchString(phone)
}

// SanitizeString sanitizes a string by removing potentially harmful characters
func SanitizeString(input string) string {
	// Remove potentially harmful characters
	re := regexp.MustCompile(`[<>'"&]`)
	return re.ReplaceAllString(input, "")
}
