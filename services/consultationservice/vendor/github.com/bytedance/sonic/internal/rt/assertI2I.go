// +build go1.17



package rt

import (
    _ `unsafe`
)

func AssertI2I(t *GoType, i GoIface) (r GoIface) {
    inter := IfaceType(t)
	tab := i.Itab
	if tab == nil {
		return
	}
	if (*GoInterfaceType)(tab.it) != inter {
		tab = GetItab(inter, tab.Vt, true)
		if tab == nil {
			return
		}
	}
	r.Itab = tab
	r.Value = i.Value
	return
}


