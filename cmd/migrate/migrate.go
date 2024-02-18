package main

import (
	"ton-lessons/internal/app"
	"ton-lessons/internal/storage"
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

	app.DB.AutoMigrate(
		&storage.Block{},
	)
	return nil
}
