package utils

import (
	"fmt"
	"github.com/google/uuid"
)

// GenerateID returns a UUID-based unique ID with optional prefix (e.g. "loan_", "inv_", etc.)
func GenerateID(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, uuid.New().String())
}