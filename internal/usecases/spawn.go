package usecases

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/internal/domain/entities"
	"github.com/steve-nzr/goff/internal/domain/factories"
	"github.com/steve-nzr/goff/internal/domain/interfaces"
	"github.com/steve-nzr/goff/internal/domain/interfaces/repositories"
	"github.com/steve-nzr/goff/internal/domain/interfaces/usecases"
	"github.com/steve-nzr/goff/internal/domain/objects"
	"github.com/steve-nzr/goff/internal/models"
)

type spawnUseCase struct {
	idgen interfaces.IdentifierGenerator

	connectionRepository    repositories.Connection
	charactersRepository    repositories.Character
	gameCharacterRepository repositories.GameCharacter
}

func (s *spawnUseCase) PreJoin(characterID customtypes.ID) *models.UseCaseResponse {
	char, err := s.charactersRepository.Get(characterID)
	if err != nil {
		logrus.Warnf("Trying to join with unknown character id %d", characterID)
		return nil
	}

	logrus.Infof("Character %s (id: %d) is joining", char.Name, char.ID)
	return &models.UseCaseResponse{
		ResponseToCaller: &objects.FPPreJoin{},
	}
}

func (s *spawnUseCase) SpawnCharacter(client *entities.NetClient, characterID customtypes.ID) *models.UseCaseResponse {
	dbChar, err := s.charactersRepository.Get(characterID)
	if err != nil {
		logrus.Error(fmt.Errorf("trying to join with unknown character id %d", characterID))
		return nil
	}

	client.CharacterID = characterID
	s.connectionRepository.UpdateCharacterID(client)

	logrus.Infof("Character %s (id: %d) is joining", dbChar.Name, dbChar.ID)

	entityID := s.idgen.Generate()
	logrus.Infof("Assigning entity id %d", entityID)

	char := factories.NewCharacter(entityID, client.ID, dbChar)
	if err = s.gameCharacterRepository.Save(char); err != nil {
		logrus.Warnf("Cannot create GameCharacter : %s", err.Error())
		return nil
	}

	logrus.Infof("Character %s (id: %d, entity id %d) joined", char.GetName(), char.GetCharacterID(), char.GetID())

	p := new(objects.FPMergePacket).Initialize(0xFF00)
	p.AddPacket(char.GetID(), objects.FPEnvironmentAllCmdID, &objects.FPEnvironmentAll{})
	p.AddPacket(char.GetID(), 0x9910, &objects.FPWorldReadInfo{Character: char})
	p.AddPacket(char.GetID(), objects.FPAddObjCmdID, &objects.FPAddObj{Character: char})

	// TODO : Spawn for others

	return &models.UseCaseResponse{
		ResponseToCaller: p,
	}
}

func NewSpawn(connectionRepository repositories.Connection,
	charactersRepository repositories.Character,
	idgen interfaces.IdentifierGenerator,
	gameCharRepo repositories.GameCharacter) usecases.Spawn {
	return &spawnUseCase{
		connectionRepository:    connectionRepository,
		charactersRepository:    charactersRepository,
		gameCharacterRepository: gameCharRepo,
		idgen:                   idgen,
	}
}
