package main

import (
	"ton-lessons/internal/app"
	scan "ton-lessons/internal/scanner"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	if err := app.InitApp(); err != nil {
		return err
	}

	scanner, err := scan.NewScanner()
	if err != nil {
		logrus.Error(err)
		return err
	}

	scanner.Listen()

	return nil
}
