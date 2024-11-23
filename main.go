package main

import (
	"github.com/Fortress-Digital/go-rest-skeleton/cmd"
	"github.com/Fortress-Digital/go-rest-skeleton/internal/log"
	"os"
)

func main() {
	logger := log.NewLogger()

	err := cmd.Execute(logger)
	logger.Error(err.Error())
	os.Exit(1)
}
