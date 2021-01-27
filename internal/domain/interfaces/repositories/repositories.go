package repositories

import (
	"net"

	"github.com/steve-nzr/goff-server/internal/domain/customtypes"
	"github.com/steve-nzr/goff-server/internal/domain/entities"
	"github.com/steve-nzr/goff-server/internal/domain/interfaces"
)

type Connection interface {
	GetByConn(conn net.Conn) (*entities.NetClient, error)
	GetByID(id customtypes.ID) (*entities.NetClient, error)
	GetByCharacterID(id customtypes.ID) (*entities.NetClient, error)
	Insert(client *entities.NetClient) error
	DeleteByConn(conn net.Conn) error
	UpdateCharacterID(client *entities.NetClient)
}

type Account interface {
	GetByName(name string) (*entities.Account, error)
	SetAuthKey(account *entities.Account, key customtypes.ID) error
}

type Character interface {
	List(account *entities.Account) ([]*entities.Character, error)
	Get(id customtypes.ID) (*entities.Character, error)
}

type GameCharacter interface {
	Get(id customtypes.ID) (interfaces.Character, error)
	Save(char interfaces.Character) error
	GetAround(id customtypes.ID) ([]interfaces.Character, error)
	GetSameMap(id customtypes.ID) ([]interfaces.Character, error)
	Delete(id customtypes.ID) error
}
