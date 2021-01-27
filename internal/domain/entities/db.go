package entities

import (
	"time"

	"github.com/golang/geo/r3"
	"github.com/steve-nzr/goff-server/internal/domain/customtypes"
)

type Account struct {
	ID             customtypes.ID `json:"id"`
	Name           string         `json:"name"`
	AuthKey        customtypes.ID `json:"auth_key"`
	AuthKeyExpires time.Time      `json:"auth_key_expires"`
}

type Character struct {
	ID        customtypes.ID `json:"id"`
	AccountID customtypes.ID `json:"account_id"`
	SkinSetID uint32         `json:"skin_set_id"`
	HairID    uint32         `json:"hair_id"`
	HairColor uint32         `json:"hair_color"`
	HeadID    uint32         `json:"head_id"`
	Location  *Location      `json:"location"`
	Name      string         `json:"name"`
	Slot      uint8          `json:"slot"`
	Gender    uint8          `json:"gender"`
	JobID     uint8          `json:"job_id"`
	Level     uint16         `json:"level"`
	Items     []*Item        `json:"items"`
	Gold      uint32         `json:"gold"`
}

type Location struct {
	MapID uint16     `json:"map_id"`
	Pos   *r3.Vector `json:"pos"`
}
