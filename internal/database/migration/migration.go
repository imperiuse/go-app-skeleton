package migration

import (
	"fmt"

	"github.com/imperiuse/go-app-skeleton/internal/database"
)

func ApplyMigrations(d *database.DB, dtos ...any) error {
	if err := d.AutoMigrate(
		dtos[:]...,
	); err != nil {
		return fmt.Errorf("AutoMigrate problem: %w", err)
	}

	return nil
}

type Tabler interface {
	TableName() string
}
