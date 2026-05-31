package main

import (
	"errors"
	"fmt"
)

var errDB = errors.New("db connection refused")

func queryUser(id int) error {
	return errDB
}

func loadProfile(id int) error {
	if err := queryUser(id); err != nil {
		return fmt.Errorf("load profile %d: %w", id, err)
	}
	return nil
}
