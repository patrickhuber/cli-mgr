package credhub

import (
	"fmt"

	"github.com/patrickhuber/cli-mgr/config"

	credhubcli "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

type CredHubStore struct {
	Name    string
	CredHub *credhubcli.CredHub
}

func NewCredHubStore(config *CredHubStoreConfig) (*CredHubStore, error) {
	options := createOptions(config)
	ch, err := credhubcli.New(config.Server, options...)
	if err != nil {
		return nil, err
	}
	return &CredHubStore{
		CredHub: ch,
		Name:    config.Name,
	}, nil
}

func createOptions(config *CredHubStoreConfig) []credhubcli.Option {
	options := []credhubcli.Option{}
	options = append(options, credhubcli.SkipTLSValidation(config.SkipTLSValidation))
	options = append(options, credhubcli.Auth(
		auth.UaaClientCredentials(
			config.ClientID,
			config.ClientSecret)))
	return options
}

func (store *CredHubStore) GetName() string {
	return store.Name
}

func (store *CredHubStore) GetByName(name string) (config.ConfigStoreData, error) {
	ch := store.CredHub
	cred, err := ch.GetLatestVersion(name)
	if err != nil {
		return config.ConfigStoreData{}, err
	}
	return config.ConfigStoreData{
		ID:    cred.Id,
		Name:  cred.Name,
		Value: cred.Value,
	}, nil
}

func (store *CredHubStore) Delete(name string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (store *CredHubStore) GetType() string {
	return "credhub"
}

func (store *CredHubStore) Put(name string, value string) (string, error) {
	return value, fmt.Errorf("not implemented")
}