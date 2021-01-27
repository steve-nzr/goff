package interfaces

import "github.com/steve-nzr/goff-server/internal/domain/customtypes"

type IdentifierGenerator interface {
	Generate() customtypes.ID
}
