package help

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/sony/sonyflake"
)

// Uuid generates a new UUID v4 string.
// Returns a string in the format "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx".
//
// Deprecated: Use Uuid7() for database primary keys, as UUIDv7 provides
// better index performance due to its time-ordered nature.
func Uuid() string {
	return uuid.New().String()
}

// Uuid7 generates a new UUID v7 string.
// UUIDv7 is time-ordered and recommended for database primary keys.
// Benefits over UUIDv4:
//   - Time-ordered: new IDs sort after old ones
//   - Index-friendly: sequential inserts reduce B-tree page splits
//   - Extractable timestamp: creation time can be derived from the ID
//
// Returns a string in the format "xxxxxxxx-xxxx-7xxx-xxxx-xxxxxxxxxxxx".
func Uuid7() string {
	id, err := uuid.NewV7()
	if err != nil {
		// Fallback to v4 if v7 generation fails (extremely rare)
		return uuid.New().String()
	}
	return id.String()
}

// MustUuid7 generates a new UUID v7 string or panics on error.
// Use this when UUID generation failure should be fatal.
func MustUuid7() string {
	return uuid.Must(uuid.NewV7()).String()
}

// Uuid7Time extracts the timestamp from a UUID v7 string.
// Returns the Unix timestamp in milliseconds and true if successful,
// or 0 and false if the UUID is not a valid v7.
func Uuid7Time(s string) (int64, bool) {
	id, err := uuid.Parse(s)
	if err != nil {
		return 0, false
	}
	if id.Version() != 7 {
		return 0, false
	}
	// UUIDv7 stores Unix milliseconds in the first 48 bits
	ts := int64(id[0])<<40 | int64(id[1])<<32 | int64(id[2])<<24 |
		int64(id[3])<<16 | int64(id[4])<<8 | int64(id[5])
	return ts, true
}

// SF is a global Sonyflake instance for generating distributed unique IDs.
// Sonyflake generates 63-bit unique IDs inspired by Twitter's Snowflake.
//
// Default Configuration:
//   - MachineID: Lower 16 bits of the private IP address
//   - StartTime: 2014-09-01 00:00:00 UTC (Sonyflake default)
//
// For containerized/distributed environments where multiple instances may
// share the same IP (e.g., Kubernetes pods), you should replace SF with
// a custom configured instance at application startup:
//
//	func init() {
//	    help.SF = sonyflake.NewSonyflake(sonyflake.Settings{
//	        MachineID: func() (uint16, error) {
//	            // Use pod name hash, environment variable, or other unique identifier
//	            id, _ := strconv.Atoi(os.Getenv("MACHINE_ID"))
//	            return uint16(id), nil
//	        },
//	    })
//	}
//
// Note: SF may be nil if initialization fails (rare, typically only when
// the machine has no valid private IP). SID() and SIDWithError() handle
// this case gracefully.
var SF = sonyflake.NewSonyflake(sonyflake.Settings{})

// SID generates a new Sonyflake ID as a string.
// Returns an empty string if ID generation fails (e.g., SF is nil or clock overflow).
//
// For production use where errors must be handled explicitly, use SIDWithError() instead.
func SID() string {
	if SF == nil {
		return ""
	}
	id, err := SF.NextID()
	if err != nil {
		return ""
	}
	return strconv.FormatUint(id, 10)
}

// SIDWithError generates a new Sonyflake ID as a string.
// Returns the ID and any error that occurred during generation.
//
// Possible errors:
//   - ErrSonyflakeNil: SF is nil (initialization failed)
//   - Clock overflow: Time exceeded Sonyflake's 174-year limit from StartTime
//   - Clock moved backwards: System time was adjusted
func SIDWithError() (string, error) {
	if SF == nil {
		return "", ErrSonyflakeNil
	}
	id, err := SF.NextID()
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(id, 10), nil
}

// ErrSonyflakeNil is returned when SF is nil (Sonyflake initialization failed).
var ErrSonyflakeNil = &sonyflakeError{"sonyflake: not initialized, SF is nil"}

type sonyflakeError struct {
	msg string
}

func (e *sonyflakeError) Error() string {
	return e.msg
}
