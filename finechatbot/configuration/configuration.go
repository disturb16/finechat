package configuration

import (
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// RabbitMQConfiguration represents the RabbitMQ configuration.
type RabbitMQConfiguration struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Configuration represents the configuration of the application.
type Configuration struct {
	RabbitMQ RabbitMQConfiguration `yaml:"RabbitMQ"`
}

const defaultFilePath = "configuration.yml"

var (
	// ErrNoFile is returned when the configuration file does not exist.
	ErrNoFile error = errors.New("file not found")
	// ErrParsingFile is returned when the configuration file is not valid.
	ErrParsingFile error = errors.New("unable to parse file")
)

// Get returns the configuration from the given file path.
func Get(filepath string) (*Configuration, error) {
	var config *Configuration
	var err error
	var confFile []byte

	confFile, err = ioutil.ReadFile(filepath)
	if err != nil {
		log.Println("Error reading configuration file:", filepath)
		return nil, ErrNoFile
	}

	//if file exists use its variables
	err = yaml.Unmarshal(confFile, &config)
	if err != nil {
		return nil, ErrParsingFile
	}
	return config, nil
}

// GetDefault returns the configuration from the default file path.
func GetDefault() (*Configuration, error) {
	return Get(defaultFilePath)
}
