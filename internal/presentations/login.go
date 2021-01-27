package presentations

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/steve-nzr/goff/internal/config"
	"github.com/steve-nzr/goff/internal/domain/interfaces/repositories"
	"github.com/steve-nzr/goff/internal/domain/interfaces/usecases"
	"github.com/steve-nzr/goff/internal/domain/objects"

	"github.com/steve-nzr/goff/internal/domain/entities"
	"github.com/steve-nzr/goff/internal/domain/interfaces"
	usecasesimpl "github.com/steve-nzr/goff/internal/usecases"
	"github.com/steve-nzr/goff/pkg/network"
)

type loginServer struct {
	clients map[net.Conn]*entities.NetClient

	accountRepository repositories.Account

	idgen interfaces.IdentifierGenerator

	welcomeUseCase usecases.Welcome
	loginUseCase   usecases.Login
}

func (l *loginServer) OnConnect(conn net.Conn) {
	client := &entities.NetClient{
		Writer: conn,
	}

	l.clients[conn] = client

	id, res := l.welcomeUseCase.Greet()
	client.ID = id
	_, _ = client.Write(res.Serialize())
}

func (l *loginServer) OnDisconnect(conn net.Conn) {
	if client, ok := l.clients[conn]; ok && client != nil {
		account, _ := l.accountRepository.GetByName(client.AccountName)
		if account != nil {
			_ = l.accountRepository.SetAuthKey(account, 0)
		}
	}

	delete(l.clients, conn)
}

func (l *loginServer) OnMessage(msg []byte, conn net.Conn) {
	/*defer func() {
		if r := recover(); r != nil {
			logrus.WithField("ip", conn.RemoteAddr().String()).Error(r)
		}
	}()*/

	client, ok := l.clients[conn]
	if !ok {
		logrus.WithField("ip", conn.RemoteAddr().String()).Errorf("Received message but client not found")
		return
	}

	fmt.Println(msg)
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

	cmd := p.ReadUInt32()
	logrus.
		WithField("ip", conn.RemoteAddr().String()).
		Infof("Received packet cmd 0x%x of len %d", cmd, length)

	if cmd == 0xfc {
		_ = p.ReadString() // version

		accountName := p.ReadString()
		fmt.Println("account :", accountName)

		// TODO : ValidateCredentials

		account, _ := l.accountRepository.GetByName(accountName)
		account.AuthKey = l.idgen.Generate()
		_ = l.accountRepository.SetAuthKey(account, account.AuthKey)
		client.AccountID = account.ID
		client.AccountName = account.Name

		if res := l.loginUseCase.ListServers(account); res != nil {
			if res.ResponseToCaller != nil {
				_, _ = client.Write(res.ResponseToCaller.Serialize())
			}
			// TODO : iterate on res.ResponsesToOthers
		}
	}
}

func NewLoginServer(accountRepository repositories.Account, idgen interfaces.IdentifierGenerator) network.NetHandlerUseCase {
	return &loginServer{
		clients:           make(map[net.Conn]*entities.NetClient),
		accountRepository: accountRepository,
		idgen:             idgen,
		welcomeUseCase:    usecasesimpl.NewWelcome(idgen),
		loginUseCase:      usecasesimpl.NewLogin(config.Servers, nil), // TODO : No constant config.Servers
	}
}
