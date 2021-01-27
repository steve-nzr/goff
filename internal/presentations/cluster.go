package presentations

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/internal/domain/interfaces/repositories"
	"github.com/steve-nzr/goff/internal/domain/interfaces/usecases"
	"github.com/steve-nzr/goff/internal/domain/objects"
	"github.com/steve-nzr/goff/internal/models"

	"github.com/steve-nzr/goff/internal/domain/entities"
	"github.com/steve-nzr/goff/internal/domain/interfaces"
	usecasesimpl "github.com/steve-nzr/goff/internal/usecases"
	"github.com/steve-nzr/goff/pkg/network"
)

type clusterServer struct {
	clients map[net.Conn]*entities.NetClient

	welcomeUseCase       usecases.Welcome
	characterListUseCase usecases.CharacterList
	spawnUseCase         usecases.Spawn

	accountRepository repositories.Account
}

func (c *clusterServer) OnConnect(conn net.Conn) {
	client := &entities.NetClient{
		Writer: conn,
	}

	c.clients[conn] = client
	id, res := c.welcomeUseCase.Greet()
	client.ID = id
	_, _ = client.Write(res.Serialize())
}

func (c *clusterServer) OnDisconnect(conn net.Conn) {
	delete(c.clients, conn)
}

func (c *clusterServer) OnMessage(msg []byte, conn net.Conn) {
	/*defer func() {
		if r := recover(); r != nil {
			logrus.
				WithField("ip", conn.RemoteAddr().String()).
				Error("Recovered an error in OnMessage : %v", r)
		}
	}()*/

	p := new(objects.FPReader).Initialize(msg)
	header := p.ReadByte()
	if header != 0x5e {
		logrus.
			WithField("ip", conn.RemoteAddr().String()).
			Infof("Received malformed packet")
		return
	}

	_ = p.ReadInt32() // checksum

	length := p.ReadInt32()

	_ = p.ReadInt32() // checksum2
	_ = p.ReadInt32() // -1 always

	cmd := p.ReadInt32()
	logrus.
		WithField("ip", conn.RemoteAddr().String()).
		Infof("Received packet cmd 0x%x of len %d", cmd, length)

	var res *models.UseCaseResponse

	if cmd == 0xf6 {
		protocolVersion := p.ReadString()
		fmt.Println("protocol version", protocolVersion)

		authKey := (customtypes.ID)(p.ReadInt32())
		fmt.Println("auth key", authKey)

		accountName := p.ReadString()
		fmt.Println("accountname", accountName)

		account, _ := c.accountRepository.GetByName(accountName)
		if authKey != account.AuthKey {
			logrus.Warnf("Account %s tried to fake authKey", accountName)
			return
		}

		res = c.characterListUseCase.List(account)

		_, _ = conn.Write(c.characterListUseCase.GetWorldAddress().Serialize())
	} else if cmd == 0xff05 {
		_ = p.ReadString() // account name
		characterID := p.ReadInt32()
		_ = p.ReadString() // characterName

		res = c.spawnUseCase.PreJoin((customtypes.ID)(characterID))
	}

	if res != nil {
		if res.ResponseToCaller != nil {
			_, _ = conn.Write(res.ResponseToCaller.Serialize())
		}
	}
}

func NewClusterServer(accountRepository repositories.Account, charactersRepository repositories.Character, idgen interfaces.IdentifierGenerator) network.NetHandlerUseCase {
	return &clusterServer{
		clients:              make(map[net.Conn]*entities.NetClient),
		accountRepository:    accountRepository,
		welcomeUseCase:       usecasesimpl.NewWelcome(idgen),
		characterListUseCase: usecasesimpl.NewCharacterList(charactersRepository),
		spawnUseCase:         usecasesimpl.NewSpawn(nil, charactersRepository, idgen, nil),
	}
}
