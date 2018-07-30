package main

import (
	"wen/svc-d/config"
)

type Check struct {
	*config.Config
}

/*
func (c Check) CheckAndStore() error {
	if c.Storage == nil {
		return fmt.Errorf("no storage mechanism defined")
	}
	results, err := c.Check()
	if err != nil {
		return err
	}

	//spew.Dump("results", results)

	for _, v := range results {
		fmt.Println(v.Title, v.Healthy, v.Endpoint, v.Status())
	}

	err = c.Storage.Store(results)
	if err != nil {
		return err
	}

	if m, ok := c.Storage.(checkup.Maintainer); ok {
		return m.Maintain()
	}

	return nil
}
*/
