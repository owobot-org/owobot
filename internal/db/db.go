/*
 * owobot - Your server's guardian and entertainer
 * Copyright (C) 2023 owobot Contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package db

import (
	"context"
	"embed"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed migrations
var migrations embed.FS

var db *sqlx.DB

// DB returns the global database instance
func DB() *sqlx.DB {
	return db
}

// Init opens the database and applies migrations
func Init(ctx context.Context, dsn string) error {
	g, err := sqlx.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	db = g
	return migrate(ctx, db)
}

func Close() error {
	return db.Close()
}

// version returns the current version of the database.
func version(ctx context.Context, db *sqlx.DB) string {
	var out string
	row := db.QueryRowxContext(ctx, "SELECT current FROM version")
	_ = row.Scan(&out)
	if out == "" {
		out = "0.sql"
	}
	return out
}

// migrate applies database migrations using the embedded sql files.
func migrate(ctx context.Context, db *sqlx.DB) error {
	current := version(ctx, db)
	return fs.WalkDir(migrations, "migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If the file is a directory, is not an sql file, or is not newer than current,
		// skip it.
		if d.IsDir() || filepath.Ext(path) != ".sql" || d.Name() <= current {
			return nil
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		// Open the sql file containing the migration code
		fl, err := migrations.Open(path)
		if err != nil {
			return err
		}
		defer fl.Close()

		// Read the file
		data, err := io.ReadAll(fl)
		if err != nil {
			return err
		}

		// Execute the migration
		_, err = tx.ExecContext(ctx, string(data))
		if err != nil {
			return err
		}

		// Update the version number
		_, err = tx.ExecContext(ctx, "DELETE FROM version; INSERT INTO version VALUES (?)", d.Name())
		if err != nil {
			return err
		}

		return tx.Commit()
	})
}
