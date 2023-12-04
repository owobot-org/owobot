/*
 * owobot - The coolest Discord bot ever written
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
	"database/sql"
	"errors"
)

func AddToStarboard(msgID string) error {
	_, err := db.Exec("INSERT OR ABORT INTO starboard VALUES (?)", msgID)
	return err
}

func ExistsInStarboard(msgID string) (bool, error) {
	var out bool
	row := db.QueryRowx("SELECT 1 FROM starboard WHERE id = ?", msgID)
	err := row.Scan(&out)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
