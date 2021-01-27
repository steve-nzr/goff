package factories

import (
	"github.com/steve-nzr/goff/internal/domain/aggregates"
	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/internal/domain/entities"
	"github.com/steve-nzr/goff/internal/domain/interfaces"
)

func NewCharacter(id, connId customtypes.ID, from *entities.Character) interfaces.Character {
	char := &aggregates.Character{
		Human: entities.Human{
			Mover: entities.Mover{
				Entity: entities.Entity{
					ID:       id,
					Location: from.Location,
					Angle:    0,
					Name:     from.Name,
				},
				Level:       from.Level,
				Destination: nil,
			},
			Gender:    from.Gender,
			SkinSetID: from.SkinSetID,
			HairID:    from.HairID,
			HairColor: from.HairColor,
			FaceID:    from.HeadID,
		},
		CharacterID:  from.ID,
		Slot:         from.Slot,
		ConnectionID: connId,
	}
	for i := range char.Inventory {
		char.Inventory[i] = &entities.Item{
			ItemID:   -1,
			UniqueID: -1,
			Position: (int16)(i),
			Count:    0,
		}
	}
	for _, item := range from.Items {
		if item.ItemID > 0 {
			char.Inventory[int(item.Position)] = item
		}
	}
	return char
}
