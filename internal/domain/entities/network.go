package entities

import (
	"github.com/steve-nzr/goff-server/internal/domain/customtypes"
	"io"
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
