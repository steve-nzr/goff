package usecases

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/steve-nzr/goff-server/internal/domain/entities"
	"github.com/steve-nzr/goff-server/internal/domain/objects"
	"github.com/steve-nzr/goff-server/internal/models"
	"github.com/steve-nzr/goff-server/pkg/testutils/mock_repositories"
	"github.com/stretchr/testify/suite"
)

type characterListSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	mockCharRepo *mock_repositories.MockCharacter
	uc           *characterList
}

func (s *characterListSuite) TestList() {
	suites := []struct {
		calls  []*gomock.Call
		input  *entities.Account
		expect *models.UseCaseResponse
	}{
		{
			calls: []*gomock.Call{
				s.mockCharRepo.EXPECT().List(gomock.Any()).Return(nil, fmt.Errorf("err")),
			},
			input:  nil,
			expect: nil,
		},
		{
			calls: []*gomock.Call{
				s.mockCharRepo.EXPECT().List(gomock.Any()).Return([]*entities.Character{
					{
						ID: 1,
					},
				}, nil),
			},
			input: &entities.Account{
				AuthKey: 1,
			},
			expect: &models.UseCaseResponse{
				ResponseToCaller: &objects.FPCharacterList{
					AuthKey: 1,
					Characters: []*entities.Character{
						{
							ID: 1,
						},
					},
				},
			},
		},
	}

	for _, ts := range suites {
		s.Equal(ts.expect, s.uc.List(ts.input))
	}
}

func (s *characterListSuite) TestCreate() {
	s.Panics(func() {
		s.uc.Create()
	})
}

func (s *characterListSuite) TestDelete() {
	s.Panics(func() {
		s.uc.Delete(nil, 1)
	})
}

func (s *characterListSuite) TestGetWorldAddress() {
	s.Equal(&objects.FPWorldAddress{
		Address: "127.0.0.1",
	}, s.uc.GetWorldAddress())
}

func TestCharacterList(t *testing.T) {
	s := &characterListSuite{
		ctrl: gomock.NewController(t),
	}
	s.mockCharRepo = mock_repositories.NewMockCharacter(s.ctrl)
	s.uc = NewCharacterList(s.mockCharRepo).(*characterList)

	suite.Run(t, s)
}
