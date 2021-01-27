package files

import (
	"encoding/json"
	"io/ioutil"

	"github.com/steve-nzr/goff/internal/config/constants"
	"github.com/steve-nzr/goff/internal/domain/customtypes"
	"github.com/steve-nzr/goff/internal/domain/entities"
	"github.com/steve-nzr/goff/internal/domain/interfaces/repositories"
)

type charactersRepository struct{}

func (c *charactersRepository) Get(id customtypes.ID) (*entities.Character, error) {
	data, err := ioutil.ReadFile("./data/characters.json")
	if err != nil {
		return nil, constants.ErrCharactersLoading
	}

	allCharacters := make([]*entities.Character, 0)
	_ = json.Unmarshal(data, &allCharacters)

	for i := range allCharacters {
		if allCharacters[i].ID == id {
			return allCharacters[i], nil
		}
	}
	return nil, constants.ErrNotFound
}

func (c *charactersRepository) List(account *entities.Account) ([]*entities.Character, error) {
	data, err := ioutil.ReadFile("./data/characters.json")
	if err != nil {
		return nil, constants.ErrCharactersLoading
	}

	allCharacters := make([]*entities.Character, 0)
	_ = json.Unmarshal(data, &allCharacters)

	accountCharacters := make([]*entities.Character, 0)
	for i := range allCharacters {
		if allCharacters[i].AccountID == account.ID {
			accountCharacters = append(accountCharacters, allCharacters[i])
		}
	}
	return accountCharacters, nil
}

func NewCharactersRepository() repositories.Character {
	return &charactersRepository{}
}
