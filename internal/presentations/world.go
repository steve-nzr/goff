package presentations

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/internal/domain/entities"
	"github.com/steve-nzr/goff/internal/domain/interfaces"
	"github.com/steve-nzr/goff/internal/domain/interfaces/repositories"
	"github.com/steve-nzr/goff/internal/domain/interfaces/usecases"
	"github.com/steve-nzr/goff/internal/domain/objects"
	"github.com/steve-nzr/goff/internal/models"
	usecasesimpl "github.com/steve-nzr/goff/internal/usecases"
	"github.com/steve-nzr/goff/pkg/network"
)

type worldServer struct {
	welcomeUseCase       usecases.Welcome
	characterListUseCase usecases.CharacterList
	spawnUseCase         usecases.Spawn

	accountRepository    repositories.Account
	connectionRepository repositories.Connection
}

func (c *worldServer) OnConnect(conn net.Conn) {
	client := &entities.NetClient{
		Writer: conn,
	}

	id, res := c.welcomeUseCase.Greet()
	client.ID = id

	if err := c.connectionRepository.Insert(client); err != nil {
		logrus.Errorf("Cannot insert new client %s", err.Error())
		_ = conn.Close()
		return
	}

	_, _ = client.Write(res.Serialize())
}

func (c *worldServer) OnDisconnect(conn net.Conn) {
	if err := c.connectionRepository.DeleteByConn(conn); err != nil {
		logrus.Errorf("Cannot delete client connection : %s", err.Error())
	}
}

func (c *worldServer) OnMessage(msg []byte, conn net.Conn) {
	/*defer func() {
		if r := recover(); r != nil {
			logrus.WithField("ip", conn.RemoteAddr().String()).Error("Recovered an error in OnMessage")
		}
	}()*/

	client, err := c.connectionRepository.GetByConn(conn)
	if err != nil {
		logrus.Errorf("Error while receiving message : %s", err.Error())
		_ = conn.Close()
		return
	}

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

	cmd := p.ReadUInt32()
	logrus.
		WithField("ip", conn.RemoteAddr().String()).
		Infof("Received packet cmd 0x%x of len %d", cmd, length)

	var res *models.UseCaseResponse
	if cmd == 0xff00 {
		_ = p.ReadInt32()
		characterID := p.ReadInt32()
		authKey := (customtypes.ID)(p.ReadInt32())
		_ = p.ReadInt32()  // pMover->m_idparty
		_ = p.ReadInt32()  // pMover->m_idGuild
		_ = p.ReadInt32()  // pMover->m_idWar
		_ = p.ReadInt32()  // uIdofMulti
		_ = p.ReadByte()   // nSlot
		_ = p.ReadString() // pMover->GetName()

		account, _ := c.accountRepository.GetByName(p.ReadString())
		if authKey != account.AuthKey {
			fmt.Println("Invalid authkey")
			return
		}

		res = c.spawnUseCase.SpawnCharacter(client, (customtypes.ID)(characterID))
	}

	if res != nil {
		if res.ResponseToCaller != nil {
			_, _ = client.Write(res.ResponseToCaller.Serialize())
		}
	}
}

func NewWorldServer(idgen interfaces.IdentifierGenerator,
	charRepo repositories.Character,
	connRepo repositories.Connection,
	accountRepo repositories.Account,
	gameCharRepo repositories.GameCharacter) network.NetHandlerUseCase {
	return &worldServer{
		welcomeUseCase:       usecasesimpl.NewWelcome(idgen),
		characterListUseCase: usecasesimpl.NewCharacterList(charRepo),
		spawnUseCase:         usecasesimpl.NewSpawn(connRepo, charRepo, idgen, gameCharRepo),
		connectionRepository: connRepo,
		accountRepository:    accountRepo,
	}
}
