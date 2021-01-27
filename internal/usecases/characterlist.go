package usecases

import (
	"github.com/steve-nzr/goff-server/internal/config"
	"github.com/steve-nzr/goff-server/internal/domain/customtypes"
	"github.com/steve-nzr/goff-server/internal/domain/entities"
	"github.com/steve-nzr/goff-server/internal/domain/interfaces/repositories"
	"github.com/steve-nzr/goff-server/internal/domain/interfaces/usecases"
	"github.com/steve-nzr/goff-server/internal/domain/objects"
	"github.com/steve-nzr/goff-server/internal/models"
	"github.com/steve-nzr/goff-server/pkg/abstract"
)

type characterList struct {
	charactersRepository repositories.Character
}

func (c *characterList) List(account *entities.Account) *models.UseCaseResponse {
	characters, err := c.charactersRepository.List(account)
	if err != nil {
		return nil
	}

	return &models.UseCaseResponse{
		ResponseToCaller: &objects.FPCharacterList{
			AuthKey:    account.AuthKey,
			Characters: characters,
		},
	}
}

func (c *characterList) Create() abstract.Serializable {
	panic("implement me")
}

func (c *characterList) Delete(account *entities.Account, characterID customtypes.ID) abstract.Serializable {
	panic("implement me")
}

func (c *characterList) GetWorldAddress() abstract.Serializable {
	return &objects.FPWorldAddress{
		Address: config.WorldAddress,
	}
}

func NewCharacterList(charactersRepository repositories.Character) usecases.CharacterList {
	return &characterList{
		charactersRepository: charactersRepository,
	}
}
