package config

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Config struct holding all config information given from yaml
type Config struct {
	BindPort int `mapstructure:"port" yaml:"port,omitempty"`

	BindHost string `mapstructure:"host" yaml:"host,omitempty"`

	Storage string `mapstructure:"storage" yaml:"storage,omitempty"`

	sync.Mutex
}

// GetAddress config
func (c *Config) GetAddress() string {
	c.Lock()
	defer c.Unlock()

	if c.BindPort == 0 {
		c.BindPort = 4444
	}
	return fmt.Sprintf("%s:%d", c.BindHost, c.BindPort)
}

// Persist current config struct to yaml file
func (c *Config) Persist() error {
	_ = c.GetAddress()

	out, err := yaml.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "yaml creation failed")
	}

	if err := ioutil.WriteFile(viper.ConfigFileUsed(), out, 0700); err != nil {
		return errors.Errorf(`Could not write to "%s" because: %s`, viper.ConfigFileUsed(), err)
	}
	return nil
}
