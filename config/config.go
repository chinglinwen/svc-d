package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sourcegraph/checkup"
)

// Config represents a configuration file.
type Config struct {
	filename string
	checkup.Checkup
}

// New creates a new Config object.
func New(filename string) *Config {
	config := Config{filename: filename}
	config.Reload()
	go config.watch()
	return &config
}

func (config *Config) Save() error {
	b, err := json.MarshalIndent(config.Checkup, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println("b", b)
	return ioutil.WriteFile(config.filename, b, 0666)
}

// Reload clears the config cache.
func (config *Config) Reload() error {
	c, err := primeCacheFromFile(config.filename)
	config.Checkup = c

	if err != nil {
		return err
	}

	return nil
}

func (config *Config) watch() {
	l := log.New(os.Stderr, "", 0)

	// Catch SIGHUP to automatically reload cache
	sighup := make(chan os.Signal, 1)
	signal.Notify(sighup, syscall.SIGHUP)

	for {
		<-sighup
		l.Println("Caught SIGHUP, reloading config...")
		config.Reload()
	}
}

func primeCacheFromFile(file string) (c checkup.Checkup, err error) {
	// File exists?
	if _, err = os.Stat(file); os.IsNotExist(err) {
		return
	}

	// Read file
	raw, e := ioutil.ReadFile(file)
	if err != nil {
		err = e
		return
	}

	// Unmarshal
	if err = json.Unmarshal(raw, &c); err != nil {
		return
	}

	return
}
