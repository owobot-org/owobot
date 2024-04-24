package builtins

import (
	"github.com/jmoiron/sqlx"
	"go.elara.ws/owobot/internal/db"
	"go.elara.ws/owobot/internal/db/sqltabler"
)

type sqlAPI struct {
	pluginName string
}

func (s sqlAPI) Exec(query string, args ...any) error {
	newQuery, err := sqltabler.Modify(query, "_owobot_plugin_", "_"+s.pluginName)
	if err != nil {
		return err
	}
	_, err = db.DB().Exec(newQuery, args...)
	return err
}

func (s sqlAPI) Query(query string, args ...any) ([]map[string]any, error) {
	newQuery, err := sqltabler.Modify(query, "_owobot_plugin_", "_"+s.pluginName)
	if err != nil {
		return nil, err
	}
	rows, err := db.DB().Queryx(newQuery, args...)
	if err != nil {
		return nil, err
	}
	return rowsToMap(rows)
}

func (s sqlAPI) QueryOne(query string, args ...any) (map[string]any, error) {
	newQuery, err := sqltabler.Modify(query, "_owobot_plugin_", "_"+s.pluginName)
	if err != nil {
		return nil, err
	}
	row := db.DB().QueryRowx(newQuery, args...)
	if err := row.Err(); err != nil {
		return nil, err
	}
	out := map[string]any{}
	return out, row.MapScan(out)
}

func rowsToMap(rows *sqlx.Rows) ([]map[string]any, error) {
	var out []map[string]any
	for rows.Next() {
		resultMap := map[string]any{}
		err := rows.MapScan(resultMap)
		if err != nil {
			return nil, err
		}
		out = append(out, resultMap)
	}

	return out, rows.Err()
}
