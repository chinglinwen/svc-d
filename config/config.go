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

// doesn't support, label, checker's name need to be unique?
// no way to delete?
//
// Config represents a configuration file.
type Config struct {
	filename string
	checkup.Checkup
	Index map[string]int
}

// project based, need to write concurrency checks
type Project struct {
	Name   string
	Region string
	checkup.Checkup
}

// New creates a new Config object.
func New(filename string) *Config {
	config := Config{filename: filename}
	config.Reload()
	go config.watch()
	return &config
}

// check if it exist first
func (c *Config) Add(name string, check checkup.Checker) error {
	for _, v := range c.Index {
		if v != 0 {
			return fmt.Errorf("already exist")
		}
	}
	i := len(c.Checkers)
	c.Checkers = append(c.Checkers, check)
	c.Index[name] = i
	return nil
}

func (c *Config) Delete(name string) error {

	/*
		for i, v := range config.Checkers {
			switch v.(type) {
			case checkup.HTTPChecker:
				a := v.(checkup.HTTPChecker)
			case checkup.TCPChecker:
				a := v.(checkup.TCPChecker)
			default:
				fmt.Errorf("no type found for add")
			}
		}
	*/
	i := c.Index[name]
	if i == 0 {
		return fmt.Errorf("not exist")
	}
	c.Checkers = append(c.Checkers[:i], c.Checkers[i+1:]...)
	delete(c.Index, "name")

}

func (config *Config) Save() error {
	b, err := json.MarshalIndent(config.Checkup, "", "  ")
	if err != nil {
		return err
	}
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
