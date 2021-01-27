package constants

import "fmt"

var (
	ErrLoginInvalidUserPassword = fmt.Errorf("invalid user or password")
	ErrLoginCapacityReached     = fmt.Errorf("server capacity reached")

	ErrCharactersLoading  = fmt.Errorf("cannot load characters")

	ErrNotFound = fmt.Errorf("not found")
)
