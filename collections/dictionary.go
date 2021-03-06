package collections

import "fmt"

type dictionary struct {
	data map[string]string
}

// ReadOnlyDictionary defines a dictionary with out set semantics
type ReadOnlyDictionary interface {
	Get(key string) (string, error)
	Lookup(key string) (string, bool)
	Keys() []string
}

// Dictionary defines a dictionary interface
type Dictionary interface {
	ReadOnlyDictionary
	Set(key, value string) error
	Unset(key string) error
}

// NewDictionary creates a new dictionary
func NewDictionary() Dictionary {
	return &dictionary{
		data: make(map[string]string),
	}
}

// NewDictionaryFromMap creates a new dictionary copying from map
func NewDictionaryFromMap(values map[string]string) Dictionary {
	dictionary := NewDictionary()
	for k, v := range values {
		dictionary.Set(k, v)
	}
	return dictionary
}

func (d *dictionary) Get(key string) (string, error) {
	if value, ok := d.data[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("unable to find key '%s' in dictionary", key)
}

func (d *dictionary) Set(key, value string) error {
	d.data[key] = value
	return nil
}

func (d *dictionary) Unset(key string) error {
	_, ok := d.data[key]
	if !ok {
		return fmt.Errorf("unable to find key '%s' in dictionary", key)
	}
	delete(d.data, key)
	return nil
}

func (d *dictionary) Lookup(key string) (string, bool) {
	value, ok := d.data[key]
	return value, ok
}

func (d *dictionary) Keys() []string {
	keys := make([]string, 0)
	for key := range d.data {
		keys = append(keys, key)
	}
	return keys
}
