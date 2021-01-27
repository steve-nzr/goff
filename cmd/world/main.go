package main

import (
	"context"
	"runtime"

	goredis "github.com/go-redis/redis/v8"
	"github.com/steve-nzr/goff/internal/domain/services"
	"github.com/steve-nzr/goff/internal/infrastructure/files"
	"github.com/steve-nzr/goff/internal/infrastructure/memory"
	"github.com/steve-nzr/goff/internal/infrastructure/redis"
	"github.com/steve-nzr/goff/internal/presentations"
	"github.com/steve-nzr/goff/pkg/network"
)

func main() {
	runtime.GOMAXPROCS(1)

	redisClient := goredis.NewClient(&goredis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	server := &network.Server{
		Network: "tcp",
		Address: "127.0.0.1:5400",
		Handler: presentations.NewWorldServer(
			services.NewIdentifierGenerator(),
			files.NewCharactersRepository(),
			memory.NewConnectionRepository(),
			files.NewAccountRepository(),
			redis.NewGameCharacterRepository(redisClient),
		),
	}

	server.Run()
}
