package network

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type NetHandlerUseCase interface {
	OnConnect(client net.Conn)
	OnDisconnect(client net.Conn)
	OnMessage(msg []byte, client net.Conn)
}

type clientMessagePair struct {
	Message []byte
	Client  net.Conn
}

// Server is a synchronous tcp/udp server
type Server struct {
	// tcp or udp
	Network string

	// e.g 127.0.0.1:8000
	Address string

	Handler NetHandlerUseCase

	clients                 []net.Conn
	connectingClientsMut    sync.Mutex
	connectingClients       []net.Conn
	disconnectingClientsMut sync.Mutex
	disconnectingClients    []net.Conn
	incomingMessagesMut     sync.Mutex
	incomingMessages        []*clientMessagePair
}

func (s *Server) Run() error {
	lis, err := net.Listen(s.Network, s.Address)
	if err != nil {
		return err
	}

	go s.handleNewClients(lis)

	for {
		s.connectNewClients()
		s.handleNewMessages()
		s.handleDisconnectingClients()

		time.Sleep(1 * time.Millisecond)
	}
}

func (s *Server) connectNewClients() {
	s.connectingClientsMut.Lock()

	for i := range s.connectingClients {
		s.Handler.OnConnect(s.connectingClients[i])
		logrus.Infof("[GOFF Network] Client connected (%s)", s.connectingClients[i].RemoteAddr().String())
	}
	s.clients = append(s.clients, s.connectingClients...)
	s.connectingClients = nil

	s.connectingClientsMut.Unlock()
}

func (s *Server) handleDisconnectingClients() {
	s.disconnectingClientsMut.Lock()

	newClients := make([]net.Conn, 0, len(s.clients))
	for i := range s.clients {
		deleting := false
		for j := range s.disconnectingClients {
			if s.clients[i] == s.disconnectingClients[j] {
				s.Handler.OnDisconnect(s.clients[i])
				logrus.Infof("[GOFF Network] Client disconnected (%s)", s.clients[i].RemoteAddr().String())
				deleting = true
				break
			}
		}
		if !deleting {
			newClients = append(newClients, s.clients[i])
		}
	}
	s.clients = newClients

	s.disconnectingClientsMut.Unlock()
}

func (s *Server) handleNewMessages() {
	s.incomingMessagesMut.Lock()

	for i := range s.incomingMessages {
		s.Handler.OnMessage(s.incomingMessages[i].Message, s.incomingMessages[i].Client)
		logrus.Infof("[GOFF Network] Handling message from client (%s)", s.incomingMessages[i].Client.RemoteAddr().String())
	}
	s.incomingMessages = nil

	s.incomingMessagesMut.Unlock()
}

func (s *Server) handleClientMessages(client net.Conn) {
	errCnt := 0
	for {
		buf := make([]byte, 4096)
		n, err := client.Read(buf)
		if err != nil {
			if errCnt >= 3 || err == io.EOF {
				logrus.Infof("[GOFF Network] Client disconnecting (%s)", client.RemoteAddr().String())
				s.disconnectingClientsMut.Lock()
				s.disconnectingClients = append(s.disconnectingClients, client)
				s.disconnectingClientsMut.Unlock()
				break
			}

			errCnt++
			fmt.Println(err)
			time.Sleep(50 * time.Millisecond)
			continue
		}

		errCnt = 0

		if n < 1 {
			fmt.Println("zero len read")
			continue
		}

		logrus.Infof("[GOFF Network] Incoming message from client (%s)", client.RemoteAddr().String())

		s.incomingMessagesMut.Lock()
		s.incomingMessages = append(s.incomingMessages, &clientMessagePair{
			Message: buf[:n],
			Client:  client,
		})
		s.incomingMessagesMut.Unlock()
	}
}

func (s *Server) handleNewClients(lis net.Listener) {
	for {
		client, err := lis.Accept()
		if err != nil {
			fmt.Println("err accept")
			continue
		}

		logrus.Infof("[GOFF Network] New client connecting (%s)", client.RemoteAddr().String())

		s.connectingClientsMut.Lock()
		s.connectingClients = append(s.connectingClients, client)
		s.connectingClientsMut.Unlock()

		go s.handleClientMessages(client)
	}
}
