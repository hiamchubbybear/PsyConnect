



package impl

import (
	"bytes"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)

func equal(in protoiface.EqualInput) protoiface.EqualOutput {
	return protoiface.EqualOutput{Equal: equalMessage(in.MessageA, in.MessageB)}
}




func equalMessage(mx, my protoreflect.Message) bool {
	if mx == nil || my == nil {
		return mx == my
	}
	if mx.Descriptor() != my.Descriptor() {
		return false
	}

	msx, ok := mx.(*messageState)
	if !ok {
		return protoreflect.ValueOfMessage(mx).Equal(protoreflect.ValueOfMessage(my))
	}
	msy, ok := my.(*messageState)
	if !ok {
		return protoreflect.ValueOfMessage(mx).Equal(protoreflect.ValueOfMessage(my))
	}

	mi := msx.messageInfo()
	miy := msy.messageInfo()
	if mi != miy {
		return protoreflect.ValueOfMessage(mx).Equal(protoreflect.ValueOfMessage(my))
	}
	mi.init()
	
	
	
	for _, ri := range mi.rangeInfos {
		var fd protoreflect.FieldDescriptor
		var vx, vy protoreflect.Value

		switch ri := ri.(type) {
		case *fieldInfo:
			hx := ri.has(msx.pointer())
			hy := ri.has(msy.pointer())
			if hx != hy {
				return false
			}
			if !hx {
				continue
			}
			fd = ri.fieldDesc
			vx = ri.get(msx.pointer())
			vy = ri.get(msy.pointer())
		case *oneofInfo:
			fnx := ri.which(msx.pointer())
			fny := ri.which(msy.pointer())
			if fnx != fny {
				return false
			}
			if fnx <= 0 {
				continue
			}
			fi := mi.fields[fnx]
			fd = fi.fieldDesc
			vx = fi.get(msx.pointer())
			vy = fi.get(msy.pointer())
		}

		if !equalValue(fd, vx, vy) {
			return false
		}
	}

	
	
	
	emx := mi.extensionMap(msx.pointer())
	emy := mi.extensionMap(msy.pointer())
	if emx != nil {
		for k, x := range *emx {
			xd := x.Type().TypeDescriptor()
			xv := x.Value()
			var y ExtensionField
			ok := false
			if emy != nil {
				y, ok = (*emy)[k]
			}
			
			if emy == nil || !ok {
				if xd.IsList() && xv.List().Len() == 0 {
					continue
				}
				return false
			}

			if !equalValue(xd, xv, y.Value()) {
				return false
			}
		}
	}
	if emy != nil {
		
		for k, y := range *emy {
			if emx != nil {
				
				if _, ok := (*emx)[k]; ok {
					continue
				}
			}
			
			if y.Type().TypeDescriptor().IsList() && y.Value().List().Len() == 0 {
				continue
			}

			
			return false
		}
	}

	return equalUnknown(mx.GetUnknown(), my.GetUnknown())
}

func equalValue(fd protoreflect.FieldDescriptor, vx, vy protoreflect.Value) bool {
	
	if fd.Kind() != protoreflect.MessageKind {
		return vx.Equal(vy)
	}

	
	if fd.IsMap() {
		if fd.MapValue().Kind() == protoreflect.MessageKind {
			return equalMessageMap(vx.Map(), vy.Map())
		}
		return vx.Equal(vy)
	}

	if fd.IsList() {
		return equalMessageList(vx.List(), vy.List())
	}

	return equalMessage(vx.Message(), vy.Message())
}




func equalMessageMap(mx, my protoreflect.Map) bool {
	if mx.Len() != my.Len() {
		return false
	}
	equal := true
	mx.Range(func(k protoreflect.MapKey, vx protoreflect.Value) bool {
		if !my.Has(k) {
			equal = false
			return false
		}
		vy := my.Get(k)
		equal = equalMessage(vx.Message(), vy.Message())
		return equal
	})
	return equal
}



func equalMessageList(lx, ly protoreflect.List) bool {
	if lx.Len() != ly.Len() {
		return false
	}
	for i := 0; i < lx.Len(); i++ {
		
		if !equalMessage(lx.Get(i).Message(), ly.Get(i).Message()) {
			return false
		}
	}
	return true
}




func equalUnknown(x, y protoreflect.RawFields) bool {
	if len(x) != len(y) {
		return false
	}
	if bytes.Equal([]byte(x), []byte(y)) {
		return true
	}

	mx := make(map[protoreflect.FieldNumber]protoreflect.RawFields)
	my := make(map[protoreflect.FieldNumber]protoreflect.RawFields)
	for len(x) > 0 {
		fnum, _, n := protowire.ConsumeField(x)
		mx[fnum] = append(mx[fnum], x[:n]...)
		x = x[n:]
	}
	for len(y) > 0 {
		fnum, _, n := protowire.ConsumeField(y)
		my[fnum] = append(my[fnum], y[:n]...)
		y = y[n:]
	}
	if len(mx) != len(my) {
		return false
	}

	for k, v1 := range mx {
		if v2, ok := my[k]; !ok || !bytes.Equal([]byte(v1), []byte(v2)) {
			return false
		}
	}

	return true
}
