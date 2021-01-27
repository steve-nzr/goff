package interfaces

import "github.com/steve-nzr/goff/internal/domain/customtypes"

type IdentifierGenerator interface {
	Generate() customtypes.ID
}
