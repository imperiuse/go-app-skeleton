package migrations

import (
	"embed"
)

// nolint
// Migrations - embed migrations
//go:embed migrations/*
var Migrations embed.FS
