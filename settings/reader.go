package settings

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Reader defines a settings reader
type Reader interface {
	Read() (*Settings, error)
}

type reader struct {
	rd io.Reader
}

// NewReader creates a new settings reader from the given io reader
func NewReader(r io.Reader) Reader {
	return &reader{
		rd: r,
	}
}

func (reader *reader) Read() (*Settings, error) {

	content, err := ioutil.ReadAll(reader.rd)
	if err != nil {
		return nil, err
	}

	settings := &Settings{}
	err = yaml.UnmarshalStrict(content, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
