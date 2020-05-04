package config

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"os"
	"path/filepath"
)

const jsonFileExtension = ".json"

type builder interface {
	getContent() ([]byte, error)
	build(content []byte) error
	validate() error
	getConfig() *Config
}


func GetConfig(path string) (*Config, error) {
	object, err := createObject(path)
	if err != nil {
		return nil, err
	}

	content, err := object.getContent()
	if err != nil {
		return nil, err
	}

	err = object.build(content)
	if err != nil {
		return nil, err
	}

	err = object.validate()
	if err != nil {
		return nil, err
	}

	return object.getConfig(), nil
}

func createObject(path string) (builder, error) {
	extension := filepath.Ext(path)

	var object builder
	switch extension {
	case jsonFileExtension:
		object = &jsonConfig{
			abstract{
				path: path,
			},
		}
	}

	if object == nil {
		return nil, errors.New("unsupported format format config")
	}

	return object, nil
}


type abstract struct {
	path   string
	config *Config
}

func (a *abstract) validate() error {
	var validate = validator.New()
	return validate.Struct(a.config)
}

func (a *abstract) getConfig() *Config {
	return a.config
}

func (a *abstract) getContent() ([]byte, error) {
	file, err := os.Open(a.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

