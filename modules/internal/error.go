package internal

import "fmt"

var (
	ErrNotFound     = fmt.Errorf("key not found")
	ErrTypeMismatch = fmt.Errorf("key type mismatch")
)
