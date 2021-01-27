package aggregates

import (
	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/internal/domain/entities"
)

type Character struct {
	entities.Human
	entities.Inventory

	CharacterID  customtypes.ID
	ConnectionID customtypes.ID

	Slot uint8
}

func (c *Character) GetCharacterID() customtypes.ID {
	return c.CharacterID
}

func (c *Character) GetGold() uint32 {
	return 0
}

func (c *Character) GetSlot() uint8 {
	return c.Slot
}

func (c *Character) GetJob() uint8 {
	return 0
}
