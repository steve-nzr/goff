package usecases

import (
	"github.com/sirupsen/logrus"
	"github.com/steve-nzr/goff-server/internal/config/constants"
	"github.com/steve-nzr/goff-server/internal/domain/entities"
	"github.com/steve-nzr/goff-server/internal/domain/interfaces/usecases"
	"github.com/steve-nzr/goff-server/internal/domain/objects"
	"github.com/steve-nzr/goff-server/internal/models"
)

type loginUseCase struct {
	accounts map[string]string // TODO : Repository
	servers  []*objects.Server // TODO : Another thing
	logger   logrus.FieldLogger
}

func (l *loginUseCase) ListServers(account *entities.Account) *models.UseCaseResponse {
	if len(l.servers) == 0 {
		l.logger.Error("Server list is empty")
		return &models.UseCaseResponse{
			ResponseToCaller: &objects.FPLoginError{
				Err: constants.ErrLoginCapacityReached,
			},
		}
	}

	return &models.UseCaseResponse{
		ResponseToCaller: &objects.FPServerList{
			Account: account.Name,
			AuthKey: account.AuthKey,
			Servers: l.servers,
		},
	}
}

func (l *loginUseCase) ValidateCredentials(account, password string) error {
	accountPassword, ok := l.accounts[account]
	if !ok {
		return constants.ErrLoginInvalidUserPassword
	}

	if password != accountPassword {
		return constants.ErrLoginInvalidUserPassword
	}

	return nil
}

func NewLogin(servers []*objects.Server, logger logrus.FieldLogger) usecases.Login {
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	return &loginUseCase{
		accounts: map[string]string{
			"test": "89d1ed22aac58f5bbea53b2fde81a946",
		},
		logger:  logger,
		servers: servers,
	}
}
