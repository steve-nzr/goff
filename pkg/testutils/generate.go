package testutils

// Generate golang
//go:generate mockgen -destination ./mock_io/io.go io Writer

// Generate logger
//go:generate mockgen -destination ./mock_logger/logger.go github.com/sirupsen/logrus FieldLogger

// Generate services
//go:generate mockgen -destination ./mock_interfaces/services.go github.com/steve-nzr/goff-server/internal/domain/interfaces IdentifierGenerator

// Generate usecases
//go:generate mockgen -destination ./mock_usecases/usecases.go github.com/steve-nzr/goff-server/internal/domain/interfaces/usecases Welcome,Login,CharacterList

// Generate repositories
//go:generate mockgen -destination ./mock_repositories/repositories.go github.com/steve-nzr/goff-server/internal/domain/interfaces/repositories Connection,Character,Account,GameCharacter
