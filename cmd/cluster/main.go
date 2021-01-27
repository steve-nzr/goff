package main

import (
	"runtime"

	"github.com/steve-nzr/goff-server/internal/domain/services"
	"github.com/steve-nzr/goff-server/internal/infrastructure/files"
	"github.com/steve-nzr/goff-server/internal/presentations"
	"github.com/steve-nzr/goff-server/pkg/network"
)

func main() {
	runtime.GOMAXPROCS(1)

	server := &network.Server{
		Network: "tcp",
		Address: "127.0.0.1:28000",
		Handler: presentations.NewClusterServer(
			files.NewAccountRepository(),
			files.NewCharactersRepository(),
			services.NewIdentifierGenerator(),
		),
	}

	server.Run()
}
