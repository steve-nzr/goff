package usecases

import (
	"fmt"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/steve-nzr/goff-server/internal/config/constants"
	"github.com/steve-nzr/goff-server/internal/domain/customtypes"
	"github.com/steve-nzr/goff-server/internal/domain/entities"
	"github.com/steve-nzr/goff-server/internal/domain/interfaces/usecases"
	"github.com/steve-nzr/goff-server/internal/domain/objects"
	"github.com/steve-nzr/goff-server/internal/models"
	"github.com/steve-nzr/goff-server/pkg/testutils/mock_interfaces"
	"github.com/steve-nzr/goff-server/pkg/testutils/mock_repositories"
	"github.com/stretchr/testify/suite"
)

type spawnSuite struct {
	suite.Suite
	ctrl          *gomock.Controller
	mockIDGen     *mock_interfaces.MockIdentifierGenerator
	mockConn      *mock_repositories.MockConnection
	mockCharacter *mock_repositories.MockCharacter
	mockGameChar  *mock_repositories.MockGameCharacter
	uc            usecases.Spawn
}

func (s *spawnSuite) TestPreJoin() {
	suites := []struct {
		calls  []*gomock.Call
		input  customtypes.ID
		expect *models.UseCaseResponse
	}{
		{
			calls: []*gomock.Call{
				s.mockCharacter.
					EXPECT().
					Get((customtypes.ID)(1)).
					Return(nil, fmt.Errorf("err")),
			},
			input:  1,
			expect: nil,
		},
		{
			calls: []*gomock.Call{
				s.mockCharacter.
					EXPECT().
					Get((customtypes.ID)(1)).
					Return(&entities.Character{ID: 1}, nil),
			},
			input: 1,
			expect: &models.UseCaseResponse{
				ResponseToCaller: &objects.FPPreJoin{},
			},
		},
	}

	for _, ts := range suites {
		s.Equal(ts.expect, s.uc.PreJoin(ts.input))
	}
}

func (s *spawnSuite) TestSpawnCharacter_Unknown() {
	s.mockCharacter.
		EXPECT().
		Get((customtypes.ID)(1)).
		Return(nil, fmt.Errorf(""))

	res := s.uc.SpawnCharacter(nil, 1)
	s.Nil(res)
}

func (s *spawnSuite) TestSpawnCharacter_CannotSave() {
	s.mockCharacter.
		EXPECT().
		Get((customtypes.ID)(1)).
		Return(&entities.Character{
			ID:       1,
			Location: &entities.Location{},
		}, nil)

	s.mockConn.
		EXPECT().
		UpdateCharacterID(&entities.NetClient{
			ID:          1,
			CharacterID: 1,
		})

	s.mockIDGen.
		EXPECT().
		Generate().
		Return((customtypes.ID)(1))

	s.mockGameChar.
		EXPECT().
		Save(gomock.Any()).
		Return(fmt.Errorf("db error"))

	res := s.uc.SpawnCharacter(&entities.NetClient{
		ID: 1,
	}, 1)
	s.Nil(res)
}

func (s *spawnSuite) TestSpawnCharacter_OK() {
	inventory := entities.Inventory{}
	for i := 0; i < constants.MaxItems; i++ {
		item := new(entities.Item)
		_ = faker.FakeData(item)
		inventory[i] = item
	}

	s.mockCharacter.
		EXPECT().
		Get((customtypes.ID)(1)).
		Return(&entities.Character{
			ID:       1,
			Location: &entities.Location{},
			Items:    inventory,
		}, nil)

	s.mockConn.
		EXPECT().
		UpdateCharacterID(&entities.NetClient{
			ID:          1,
			CharacterID: 1,
		})

	s.mockIDGen.
		EXPECT().
		Generate().
		Return((customtypes.ID)(1))

	s.mockGameChar.
		EXPECT().
		Save(gomock.Any()).
		Return(nil)

	res := s.uc.SpawnCharacter(&entities.NetClient{
		ID: 1,
	}, 1)
	s.NotNil(res.ResponseToCaller)
	s.Len(res.ResponsesToOthers, 0)
}

/*
func (s *spawnSuite) TestSpawnCharacter() {
	items := make([]*entities.Item, 0, constants.MaxItems)
	for i := 0; i < constants.MaxItems; i++ {
		item := new(entities.Item)
		_ = faker.FakeData(item)
		items = append(items, item)
	}

	suites := []struct {
		calls            []*gomock.Call
		inputClient      *entities.NetClient
		inputCharacterID customtypes.ID
		expect           error
	}{
		{
			calls: []*gomock.Call{
				s.mockCharacter.
					EXPECT().
					Get((customtypes.ID)(1)).
					Return(nil, fmt.Errorf("")),
			},
			inputClient:      nil,
			inputCharacterID: 1,
			expect:           fmt.Errorf(""),
		},
		{
			calls: []*gomock.Call{
				s.mockCharacter.
					EXPECT().
					Get((customtypes.ID)(1)).
					Return(&entities.Character{
						ID:       1,
						Location: &entities.Location{},
						Items:    items,
					}, nil),
				s.mockConn.
					EXPECT().
					UpdateCharacterID(&entities.NetClient{
						Writer:      s.mockWriter,
						ID:          1,
						CharacterID: 1,
					}),
				s.mockIDGen.
					EXPECT().
					Generate().
					Return((customtypes.ID)(1)),
				s.mockWriter.
					EXPECT().
					Write(gomock.Any()).
					Return(0, nil),
			},
			inputClient: &entities.NetClient{
				Writer: s.mockWriter,
				ID:     1,
			},
			inputCharacterID: 1,
			expect:           nil,
		},
	}

	for _, ts := range suites {
		err := s.uc.SpawnCharacter(ts.inputClient, ts.inputCharacterID)
		if ts.expect != nil {
			s.NotNil(err)
		} else {
			s.Nil(err)
		}
	}
}*/

func TestSpawnUseCase(t *testing.T) {
	s := &spawnSuite{
		ctrl: gomock.NewController(t),
	}
	s.mockConn = mock_repositories.NewMockConnection(s.ctrl)
	s.mockCharacter = mock_repositories.NewMockCharacter(s.ctrl)
	s.mockIDGen = mock_interfaces.NewMockIdentifierGenerator(s.ctrl)
	s.mockGameChar = mock_repositories.NewMockGameCharacter(s.ctrl)
	s.uc = NewSpawn(s.mockConn, s.mockCharacter, s.mockIDGen, s.mockGameChar)
	suite.Run(t, s)
}
