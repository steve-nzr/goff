package entities

import (
	"github.com/golang/geo/r3"
	"github.com/steve-nzr/goff-server/internal/domain/customtypes"
)

type Entity struct {
	ID       customtypes.ID
	Location *Location
	Angle    float32
	Name     string
}

func (l *Entity) GetID() customtypes.ID {
	return l.ID
}

func (l *Entity) GetType() uint8 {
	return 0
}

func (l *Entity) GetAngle() float32 {
	return l.Angle
}

func (l *Entity) GetMapID() uint16 {
	return l.Location.MapID
}

func (l *Entity) GetLocation() *r3.Vector {
	return l.Location.Pos
}

func (l *Entity) GetName() string {
	return l.Name
}

func (l *Entity) GetSize() uint16 {
	return 100
}

type Mover struct {
	Entity

	Level       uint16
	Destination *r3.Vector
}

func (c *Mover) GetLevel() uint16 {
	return c.Level
}

func (c *Mover) GetDestination(_ *r3.Vector) *r3.Vector {
	return c.Destination
}

func (c *Mover) GetSpeed() float32 {
	return 100.0
}

func (c *Mover) Move(pos *r3.Vector) error {
	c.Location.Pos = pos
	return nil
}

func (c *Mover) Teleport(mapID uint16, pos *r3.Vector) {
	panic("not implemented") // TODO: Implement
}

type Human struct {
	Mover

	Gender    uint8
	SkinSetID uint32
	HairID    uint32
	HairColor uint32
	FaceID    uint32
}

func (h *Human) GetGender() uint8 {
	return h.Gender
}

func (h *Human) GetSkinSetID() uint32 {
	return h.SkinSetID
}

func (h *Human) GetHairID() uint32 {
	return h.HairID
}

func (h *Human) GetHairColor() uint32 {
	return h.HairColor
}

func (h *Human) GetFaceID() uint32 {
	return h.FaceID
}
