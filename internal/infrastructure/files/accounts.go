package files

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/steve-nzr/goff-server/internal/config/constants"
	"github.com/steve-nzr/goff-server/internal/domain/customtypes"
	"github.com/steve-nzr/goff-server/internal/domain/entities"
	"github.com/steve-nzr/goff-server/internal/domain/interfaces/repositories"
)

type accountRepository struct{}

func (a *accountRepository) GetByName(name string) (*entities.Account, error) {
	data, err := ioutil.ReadFile("./data/accounts.json")
	if err != nil {
		return nil, constants.ErrNotFound
	}

	accounts := make([]*entities.Account, 0)
	_ = json.Unmarshal(data, &accounts)
	for i := range accounts {
		if accounts[i].Name == name {
			return accounts[i], nil
		}
	}

	return nil, constants.ErrNotFound
}

func (a *accountRepository) SetAuthKey(account *entities.Account, key customtypes.ID) error {
	accounts, err := a.readAccounts()
	if err != nil {
		return err
	}

	for i := range accounts {
		if accounts[i].ID == account.ID {
			accounts[i].AuthKey = key
			accounts[i].AuthKeyExpires = time.Now().Add(15 * time.Minute)
			return a.saveAccounts(accounts)
		}
	}

	return constants.ErrNotFound
}

func (a *accountRepository) readAccounts() ([]*entities.Account, error) {
	data, err := ioutil.ReadFile("./data/accounts.json")
	if err != nil {
		return nil, constants.ErrNotFound
	}

	accounts := make([]*entities.Account, 0)
	return accounts, json.Unmarshal(data, &accounts)
}

func (a *accountRepository) saveAccounts(accounts []*entities.Account) error {
	data, _ := json.Marshal(accounts)
	return ioutil.WriteFile("./data/accounts.json", data, 0644)
}

func NewAccountRepository() repositories.Account {
	return &accountRepository{}
}
