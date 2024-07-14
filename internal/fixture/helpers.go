// Package fixture - fixture for dev mode.
// //nolint: gomnd,gosec,unused, - it's ok here, helper package for fixtures.
package fixture

import (
	"github.com/imperiuse/go-app-skeleton/internal/database"
)

// TruncateTables - truncate DTO's tables.
func TruncateTables(db *database.DB, dtos ...any) error {
	if err := db.Model(dtos).Delete("", "1=1").Error; err != nil {
		return err
	}

	return nil
}

// FillTables - fill DTO's tables. //nolint: gosec.
func FillTables[T any](db *database.DB, dtos ...T) error {
	if err := db.Create(dtos).Error; err != nil { // //nolint: gosec // this is ok.
		return err
	}

	return nil
}
