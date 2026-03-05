

package loader

import (
    _ `unsafe`
)

//go:linkname lastmoduledatap runtime.lastmoduledatap

var lastmoduledatap *moduledata

//go:linkname moduledataverify1 runtime.moduledataverify1
func moduledataverify1(_ *moduledata)
