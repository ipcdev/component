package errors

import (
	"fmt"
)

const (
	ConfigurationNotValid int = iota + 1000
	ErrInvalidJSON
	ErrEOF
	ErrLoadConfigFailed
)

func init() {
	Register(defaultCoder{500, ConfigurationNotValid, "ConfigurationNotValid error"})
	Register(defaultCoder{500, ErrInvalidJSON, "Data is not valid JSON"})
	Register(defaultCoder{500, ErrEOF, "End of input"})
	Register(defaultCoder{500, ErrLoadConfigFailed, "Load configuration file failed"})
}

func loadConfig() error {
	err := decodeConfig()
	return Wrapc(err, ConfigurationNotValid, "service configuration could not be loaded")
}

func decodeConfig() error {
	err := readConfig()
	return Wrapc(err, ErrInvalidJSON, "could not decode configuration data")
}

func readConfig() error {
	err := fmt.Errorf("read: end of input")
	return Wrapc(err, ErrEOF, "could not read configuration file")
}
