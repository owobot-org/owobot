package builtins

import (
	"errors"

	"github.com/dop251/goja"
)

// Register registers all the owobot APIs in JavaScript.
func Register(vm *goja.Runtime, pluginName, pluginVersion string) error {
	return errors.Join(
		vm.GlobalObject().Set("sql", sqlAPI{pluginName: pluginName}),
		vm.GlobalObject().Set("vercmp", vercmpAPI{}),
		vm.GlobalObject().Set("cache", cacheAPI{}),
		vm.GlobalObject().Set("tickets", ticketsAPI{}),
		vm.GlobalObject().Set("eventlog", eventLogAPI{}),
		vm.GlobalObject().Set("fetch", fetch(pluginName, pluginVersion)),
	)
}
