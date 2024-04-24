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

type PluginInfo struct {
	Name    string `db:"name"`
	Version string `db:"version"`
	Desc    string `db:"description"`
}

func (pi PluginInfo) IsValid() bool {
	if pi.Name == "" || pi.Version == "" || pi.Desc == "" {
		return false
	}
	return true
}

func AddPlugin(pi PluginInfo) error {
	_, err := db.NamedExec(`INSERT OR REPLACE INTO plugins VALUES (:name, :version, :description)`, pi)
	return err
}

func GetPlugin(name string) (out PluginInfo, err error) {
	err = db.QueryRowx("SELECT * FROM plugins WHERE name = ? LIMIT 1", name).StructScan(&out)
	return
}
