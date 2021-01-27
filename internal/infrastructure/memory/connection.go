package memory

import (
	"io"
	"net"

	"github.com/steve-nzr/goff-server/internal/config/constants"
	"github.com/steve-nzr/goff-server/internal/domain/customtypes"
	"github.com/steve-nzr/goff-server/internal/domain/entities"
	"github.com/steve-nzr/goff-server/internal/domain/interfaces/repositories"
)

type connectionRepository struct {
	clientsByConn        map[io.Writer]*entities.NetClient
	clientsByID          map[customtypes.ID]*entities.NetClient
	clientsByCharacterID map[customtypes.ID]*entities.NetClient
}

func (c *connectionRepository) UpdateCharacterID(client *entities.NetClient) {
	delete(c.clientsByCharacterID, client.PreviousCharacterID)
	if client.CharacterID != 0 {
		c.clientsByCharacterID[client.CharacterID] = client
	}
}

func (c *connectionRepository) GetByCharacterID(id customtypes.ID) (*entities.NetClient, error) {
	client, found := c.clientsByCharacterID[id]
	if !found {
		return nil, constants.ErrNotFound
	}

	if client.CharacterID != id {
		return nil, constants.ErrNotFound
	}

	return client, nil
}

func (c *connectionRepository) GetByConn(conn net.Conn) (*entities.NetClient, error) {
	client, found := c.clientsByConn[conn]
	if !found {
		return nil, constants.ErrNotFound
	}

	return client, nil
}

func (c *connectionRepository) GetByID(id customtypes.ID) (*entities.NetClient, error) {
	client, found := c.clientsByID[id]
	if !found {
		return nil, constants.ErrNotFound
	}

	return client, nil
}

func (c *connectionRepository) Insert(client *entities.NetClient) error {
	c.clientsByConn[client.Writer] = client
	c.clientsByID[client.ID] = client
	return nil
}

func (c *connectionRepository) DeleteByConn(conn net.Conn) error {
	client, _ := c.GetByConn(conn)
	if client != nil {
		delete(c.clientsByID, client.ID)
	}
	delete(c.clientsByConn, conn)
	return nil
}

func NewConnectionRepository() repositories.Connection {
	return &connectionRepository{
		clientsByConn:        make(map[io.Writer]*entities.NetClient),
		clientsByID:          make(map[customtypes.ID]*entities.NetClient),
		clientsByCharacterID: make(map[customtypes.ID]*entities.NetClient),
	}
}
