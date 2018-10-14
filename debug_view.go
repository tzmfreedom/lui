package main

import (
	"os"

	"github.com/k0kubun/pp"
)

func debug(obj interface{}) error {
	f, err := os.Create("./debug")
	if err != nil {
		return err
	}
	defer f.Close()
	pp.Fprintln(f, obj)
	return nil
}
