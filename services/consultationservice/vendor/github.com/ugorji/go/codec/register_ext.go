


package codec

import "reflect"





func (h *JsonHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}


func (h *CborHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}


func (h *MsgpackHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}


func (h *SimpleHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}


func (h *BincHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) {
	return h.SetExt(rt, tag, makeExt(ext))
}




