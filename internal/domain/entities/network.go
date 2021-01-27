package entities

import (
	"io"

	"github.com/steve-nzr/goff/internal/domain/customtypes"
)

type NetClient struct {
	io.Writer
	ID          customtypes.ID
	CharacterID customtypes.ID
	AccountID   customtypes.ID
	AccountName string

	PreviousCharacterID customtypes.ID
}

func (c *NetClient) Write(data []byte) (n int, err error) {
	return c.Writer.Write(data)
}
