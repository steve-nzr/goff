package usecases

import (
	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/internal/domain/entities"
	"github.com/steve-nzr/goff/internal/models"
	"github.com/steve-nzr/goff/pkg/abstract"
)

type Welcome interface {
	Greet() (customtypes.ID, abstract.Serializable)
}

type Login interface {
	ValidateCredentials(account, password string) error
	ListServers(account *entities.Account) *models.UseCaseResponse
}

type CharacterList interface {
	List(account *entities.Account) *models.UseCaseResponse
	Create() abstract.Serializable
	Delete(account *entities.Account, characterID customtypes.ID) abstract.Serializable
	GetWorldAddress() abstract.Serializable
}

type Spawn interface {
	PreJoin(characterID customtypes.ID) *models.UseCaseResponse
	SpawnCharacter(client *entities.NetClient, characterID customtypes.ID) *models.UseCaseResponse
}

/*type LivingBehviourUseCase interface {
	MoveByClick(id customtypes.ID, destination r3.Vector) *models.UseCaseResponse
	DoMotion(id customtypes.ID, motionID int) *models.UseCaseResponse
}*/
