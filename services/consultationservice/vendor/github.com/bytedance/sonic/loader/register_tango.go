// +build bytedance_tango



package loader

import (
    "sync"
	_ "unsafe"
)

//go:linkname pluginsMu plugin.pluginsMu
var pluginsMu sync.Mutex

func registerModule(mod *moduledata) {
    pluginsMu.Lock()
    defer pluginsMu.Unlock()
    lastmoduledatap.next = mod
    lastmoduledatap = mod
}
