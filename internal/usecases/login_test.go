package usecases

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/steve-nzr/goff/internal/config"
	"github.com/steve-nzr/goff/internal/config/constants"
	"github.com/steve-nzr/goff/internal/domain/entities"
	"github.com/steve-nzr/goff/internal/domain/objects"
	"github.com/steve-nzr/goff/internal/models"
	mock_logrus "github.com/steve-nzr/goff/pkg/testutils/mock_logger"
	"github.com/stretchr/testify/suite"
)

type loginSuite struct {
	suite.Suite
	ctrl    *gomock.Controller
	mockLog *mock_logrus.MockFieldLogger
	client  *entities.NetClient
	uc      *loginUseCase
}

func TestLoginSuite(t *testing.T) {
	s := &loginSuite{
		ctrl: gomock.NewController(t),
	}
	s.mockLog = mock_logrus.NewMockFieldLogger(s.ctrl)

	_ = NewLogin(nil, nil)
	s.uc = NewLogin(config.Servers, s.mockLog).(*loginUseCase)

	suite.Run(t, s)
}

func (s *loginSuite) TestListServers() {
	res := s.uc.ListServers(&entities.Account{Name: "test", AuthKey: 500})

	s.Equal(&models.UseCaseResponse{
		ResponseToCaller: &objects.FPServerList{
			Account: "test",
			AuthKey: 500,
			Servers: config.Servers,
		},
	}, res)
}

func (s *loginSuite) TestListServersEmptyServers() {
	s.uc = NewLogin(nil, s.mockLog).(*loginUseCase)

	// expect
	s.mockLog.
		EXPECT().
		Error("Server list is empty")

	res := s.uc.ListServers(&entities.Account{})

	s.Equal(&models.UseCaseResponse{
		ResponseToCaller: &objects.FPLoginError{
			Err: constants.ErrLoginCapacityReached,
		},
	}, res)

	// put back usecase
	s.uc = NewLogin(config.Servers, s.mockLog).(*loginUseCase)
}

func (s *loginSuite) TestValidateCredentials() {
	suites := []struct {
		account  string
		password string
		expect   error
	}{
		{
			account:  "",
			password: "",
			expect:   constants.ErrLoginInvalidUserPassword,
		},
		{
			account:  "unknown",
			password: "",
			expect:   constants.ErrLoginInvalidUserPassword,
		},
		{
			account:  "test",
			password: "invalid",
			expect:   constants.ErrLoginInvalidUserPassword,
		},
		{
			account:  "test",
			password: "89d1ed22aac58f5bbea53b2fde81a946",
			expect:   nil,
		},
	}

	for _, ts := range suites {
		s.Equal(ts.expect, s.uc.ValidateCredentials(ts.account, ts.password))
	}
}
