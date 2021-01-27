package services

import (
	"github.com/steve-nzr/goff-server/internal/domain/interfaces"
	"github.com/stretchr/testify/suite"
	"testing"
)

type identifierGeneratorSuite struct {
	suite.Suite
	svc interfaces.IdentifierGenerator
}

func (s *identifierGeneratorSuite) TestGenerate() {
	id := s.svc.Generate()
	s.NotZero(id)
}

func TestIdentifierGenerator(t *testing.T) {
	suite.Run(t, &identifierGeneratorSuite{
		svc: NewIdentifierGenerator(),
	})
}
