package internal

import (
	"fmt"
	"strconv"
)

// ParseUint64 parses s as a decimal uint64 value.
func ParseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// ParseUint32 parses s as a decimal uint32 value.
func ParseUint32(s string) (uint32, error) {
	// In response to https://github.com/PeggyJV/sommelier/issues/292
	// let's use strconv.ParseUint to properly catch range errors and
	// avoid underflows.
	u, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	// Bullet-proof check to ensure no underflows (even though we already have range checks)
	u32 := uint32(u)
	if g := uint64(u32); g != u {
		return 0, fmt.Errorf("parseuint32 underflow detected: got %d, want %d", g, u)
	}
	return u32, nil
}
