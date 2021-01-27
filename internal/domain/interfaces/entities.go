package interfaces

import (
	"github.com/golang/geo/r3"
	"github.com/steve-nzr/goff-server/internal/domain/customtypes"
)

type Entity interface {
	GetID() customtypes.ID
	GetType() uint8

	GetMapID() uint16
	GetLocation() *r3.Vector
	GetAngle() float32

	GetName() string
	GetSize() uint16
}

type Mover interface {
	Entity

	GetLevel() uint16

	GetDestination(*r3.Vector) *r3.Vector
	GetSpeed() float32

	Move(pos *r3.Vector) error
	Teleport(mapID uint16, pos *r3.Vector)
}

type Human interface {
	Mover

	GetGender() uint8
	GetSkinSetID() uint32
	GetHairID() uint32
	GetHairColor() uint32
	GetFaceID() uint32
}

type Character interface {
	Human
	Inventory

	GetCharacterID() customtypes.ID
	GetGold() uint32
	GetSlot() uint8
	GetJob() uint8
}

type Item interface {
	GetUniqueID() int16
	GetItemID() int32
	GetPosition() int16
	GetCount() uint16
}

type Inventory interface {
	GetItemByPosition(pos int16) Item
}
