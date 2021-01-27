package entities

import (
	"github.com/steve-nzr/goff-server/internal/config/constants"
	"github.com/steve-nzr/goff-server/internal/domain/interfaces"
)

type Item struct {
	ItemID   int32  `json:"item_id" faker:"boundary_start=0, boundary_end=50000"`
	UniqueID int16  `json:"unique_id" faker:"boundary_start=0, boundary_end=73"`
	Position int16  `json:"position" faker:"boundary_start=0, boundary_end=73"`
	Count    uint16 `json:"count" faker:"oneof: 1, 1"`
}

func (i *Item) GetUniqueID() int16 {
	return i.UniqueID
}

func (i *Item) GetItemID() int32 {
	return i.ItemID
}

func (i *Item) GetPosition() int16 {
	return i.Position
}

func (i *Item) GetCount() uint16 {
	return i.Count
}

type Inventory [constants.MaxItems]*Item

func (i *Inventory) GetItemByPosition(pos int16) interfaces.Item {
	return i[pos]
}

/*func (i *Inventory) MoveItem(source, dest int) error {
	return nil
}*/
