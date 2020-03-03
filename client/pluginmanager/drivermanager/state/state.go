package state

import pstructs "github.com/actiontech/dtle/plugins/shared/structs"

// PluginState is used to store the driver managers state across restarts of the
// agent
type PluginState struct {
	// ReattachConfigs are the set of reattach configs for plugin's launched by
	// the driver manager
	ReattachConfigs map[string]*pstructs.ReattachConfig
}
