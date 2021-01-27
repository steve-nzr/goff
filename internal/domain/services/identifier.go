package services

import (
	"math/rand"
	"time"

	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/internal/domain/interfaces"
)

type identifierGenerator struct {
	src rand.Source
	gen *rand.Rand
}

func (i *identifierGenerator) Generate() customtypes.ID {
	return (customtypes.ID)(i.gen.Int31())
}

func NewIdentifierGenerator() interfaces.IdentifierGenerator {
	src := rand.NewSource(time.Now().UnixNano())
	return &identifierGenerator{
		src: src,
		gen: rand.New(src), //nolint:gosec
	}
}
