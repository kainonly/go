package help

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/sony/sonyflake"
)

// Uuid generates a new UUID v4 string.
// Returns a string in the format "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx".
func Uuid() string {
	return uuid.New().String()
}

// SF is a global Sonyflake instance for generating distributed unique IDs.
// Sonyflake generates 63-bit unique IDs inspired by Twitter's Snowflake.
var SF = sonyflake.NewSonyflake(sonyflake.Settings{})

// SID generates a new Sonyflake ID as a string.
// Returns an empty string if ID generation fails.
// Note: In high-availability scenarios, consider using SIDWithError instead.
func SID() string {
	id, err := SF.NextID()
	if err != nil {
		return ""
	}
	return strconv.FormatUint(id, 10)
}

// SIDWithError generates a new Sonyflake ID as a string.
// Returns the ID and any error that occurred during generation.
func SIDWithError() (string, error) {
	id, err := SF.NextID()
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(id, 10), nil
}
