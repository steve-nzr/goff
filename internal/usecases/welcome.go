package usecases

import (
	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/internal/domain/interfaces"
	"github.com/steve-nzr/goff/internal/domain/interfaces/usecases"
	"github.com/steve-nzr/goff/internal/domain/objects"
	"github.com/steve-nzr/goff/pkg/abstract"
)

type welcomeUseCase struct {
	idgen interfaces.IdentifierGenerator
}

func (w *welcomeUseCase) Greet() (customtypes.ID, abstract.Serializable) {
	id := w.idgen.Generate()
	return id, &objects.FPWelcome{
		ID: id,
	}
}

func NewWelcome(identifierGenerator interfaces.IdentifierGenerator) usecases.Welcome {
	return &welcomeUseCase{
		idgen: identifierGenerator,
	}
}
