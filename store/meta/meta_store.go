package meta

import (
	"fmt"
	"path/filepath"

	"github.com/patrickhuber/wrangle/store"
)

type metaStore struct {
	configFilePath   string
	configFileFolder string
	name             string
}

const (
	// ConfigFilePathKey is the default config file path key in the meta store
	ConfigFilePathKey = "config_file_path"

	// ConfigFileFolderKey is the directory of the default config file
	ConfigFileFolderKey = "config_file_folder"
)

// NewMetaStore creates a new meta store
func NewMetaStore(name, configFilePath string) store.Store {
	return &metaStore{
		configFilePath:   configFilePath,
		configFileFolder: filepath.Dir(configFilePath),
		name:             name,
	}
}

func (s *metaStore) Name() string { return s.name }

func (s *metaStore) Type() string { return "meta" }

func (s *metaStore) Put(key string, value string) (string, error) {
	return "", fmt.Errorf("meta.Put is not implemented")
}

func (s *metaStore) GetByName(name string) (store.Data, error) {
	var value string
	switch name {
	case ConfigFilePathKey:
		value = s.configFilePath
	case ConfigFileFolderKey:
		value = s.configFileFolder
	default:
		return nil, fmt.Errorf("unable to find key '%s' in meta store", name)
	}
	return store.NewData(name, name, value), nil
}
func (s *metaStore) Delete(key string) (int, error) { return 0, nil }