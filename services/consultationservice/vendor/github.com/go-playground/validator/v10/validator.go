package validator

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)


type validate struct {
	v              *Validate
	top            reflect.Value
	ns             []byte
	actualNs       []byte
	errs           ValidationErrors
	includeExclude map[string]struct{} 
	ffn            FilterFunc
	slflParent     reflect.Value 
	slCurrent      reflect.Value 
	flField        reflect.Value 
	cf             *cField       
	ct             *cTag         
	misc           []byte        
	str1           string        
	str2           string        
	fldIsPointer   bool          
	isPartial      bool
	hasExcludes    bool
}


func (v *validate) validateStruct(ctx context.Context, parent reflect.Value, current reflect.Value, typ reflect.Type, ns []byte, structNs []byte, ct *cTag) {

	cs, ok := v.v.structCache.Get(typ)
	if !ok {
		cs = v.v.extractStructCache(current, typ.Name())
	}

	if len(ns) == 0 && len(cs.name) != 0 {

		ns = append(ns, cs.name...)
		ns = append(ns, '.')

		structNs = append(structNs, cs.name...)
		structNs = append(structNs, '.')
	}

	
	
	if ct == nil || ct.typeof != typeStructOnly {

		var f *cField

		for i := 0; i < len(cs.fields); i++ {

			f = cs.fields[i]

			if v.isPartial {

				if v.ffn != nil {
					
					if v.ffn(append(structNs, f.name...)) {
						continue
					}

				} else {
					
					_, ok = v.includeExclude[string(append(structNs, f.name...))]

					if (ok && v.hasExcludes) || (!ok && !v.hasExcludes) {
						continue
					}
				}
			}

			v.traverseField(ctx, current, current.Field(f.idx), ns, structNs, f, f.cTags)
		}
	}

	
	
	
	if cs.fn != nil {

		v.slflParent = parent
		v.slCurrent = current
		v.ns = ns
		v.actualNs = structNs

		cs.fn(ctx, v)
	}
}


func (v *validate) traverseField(ctx context.Context, parent reflect.Value, current reflect.Value, ns []byte, structNs []byte, cf *cField, ct *cTag) {
	var typ reflect.Type
	var kind reflect.Kind

	current, kind, v.fldIsPointer = v.extractTypeInternal(current, false)

	var isNestedStruct bool

	switch kind {
	case reflect.Ptr, reflect.Interface, reflect.Invalid:

		if ct == nil {
			return
		}

		if ct.typeof == typeOmitEmpty || ct.typeof == typeIsDefault {
			return
		}

		if ct.typeof == typeOmitNil && (kind != reflect.Invalid && current.IsNil()) {
			return
		}

		if ct.typeof == typeOmitZero {
			return
		}

		if ct.hasTag {
			if kind == reflect.Invalid {
				v.str1 = string(append(ns, cf.altName...))
				if v.v.hasTagNameFunc {
					v.str2 = string(append(structNs, cf.name...))
				} else {
					v.str2 = v.str1
				}
				v.errs = append(v.errs,
					&fieldError{
						v:              v.v,
						tag:            ct.aliasTag,
						actualTag:      ct.tag,
						ns:             v.str1,
						structNs:       v.str2,
						fieldLen:       uint8(len(cf.altName)),
						structfieldLen: uint8(len(cf.name)),
						param:          ct.param,
						kind:           kind,
					},
				)
				return
			}

			v.str1 = string(append(ns, cf.altName...))
			if v.v.hasTagNameFunc {
				v.str2 = string(append(structNs, cf.name...))
			} else {
				v.str2 = v.str1
			}
			if !ct.runValidationWhenNil {
				v.errs = append(v.errs,
					&fieldError{
						v:              v.v,
						tag:            ct.aliasTag,
						actualTag:      ct.tag,
						ns:             v.str1,
						structNs:       v.str2,
						fieldLen:       uint8(len(cf.altName)),
						structfieldLen: uint8(len(cf.name)),
						value:          getValue(current),
						param:          ct.param,
						kind:           kind,
						typ:            current.Type(),
					},
				)
				return
			}
		}

		if kind == reflect.Invalid {
			return
		}

	case reflect.Struct:
		isNestedStruct = !current.Type().ConvertibleTo(timeType)
		
		
		
		
		
		if isNestedStruct && !v.v.requiredStructEnabled && ct != nil && ct.tag == requiredTag {
			ct = ct.next
		}
	}

	typ = current.Type()

OUTER:
	for {
		if ct == nil || !ct.hasTag || (isNestedStruct && len(cf.name) == 0) {
			
			if isNestedStruct {
				
				
				
				
				if len(cf.name) > 0 {
					ns = append(append(ns, cf.altName...), '.')
					structNs = append(append(structNs, cf.name...), '.')
				}

				v.validateStruct(ctx, parent, current, typ, ns, structNs, ct)
			}
			return
		}

		switch ct.typeof {
		case typeNoStructLevel:
			return

		case typeStructOnly:
			if isNestedStruct {
				
				
				
				
				if len(cf.name) > 0 {
					ns = append(append(ns, cf.altName...), '.')
					structNs = append(append(structNs, cf.name...), '.')
				}

				v.validateStruct(ctx, parent, current, typ, ns, structNs, ct)
			}
			return

		case typeOmitEmpty:

			
			v.slflParent = parent
			v.flField = current
			v.cf = cf
			v.ct = ct

			if !hasValue(v) {
				return
			}

			ct = ct.next
			continue

		case typeOmitZero:
			v.slflParent = parent
			v.flField = current
			v.cf = cf
			v.ct = ct

			if !hasNotZeroValue(v) {
				return
			}

			ct = ct.next
			continue

		case typeOmitNil:
			v.slflParent = parent
			v.flField = current
			v.cf = cf
			v.ct = ct

			switch field := v.Field(); field.Kind() {
			case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
				if field.IsNil() {
					return
				}
			default:
				if v.fldIsPointer && field.Interface() == nil {
					return
				}
			}

			ct = ct.next
			continue

		case typeEndKeys:
			return

		case typeDive:

			ct = ct.next

			
			
			switch kind {
			case reflect.Slice, reflect.Array:

				var i64 int64
				reusableCF := &cField{}

				for i := 0; i < current.Len(); i++ {

					i64 = int64(i)

					v.misc = append(v.misc[0:0], cf.name...)
					v.misc = append(v.misc, '[')
					v.misc = strconv.AppendInt(v.misc, i64, 10)
					v.misc = append(v.misc, ']')

					reusableCF.name = string(v.misc)

					if cf.namesEqual {
						reusableCF.altName = reusableCF.name
					} else {

						v.misc = append(v.misc[0:0], cf.altName...)
						v.misc = append(v.misc, '[')
						v.misc = strconv.AppendInt(v.misc, i64, 10)
						v.misc = append(v.misc, ']')

						reusableCF.altName = string(v.misc)
					}
					v.traverseField(ctx, parent, current.Index(i), ns, structNs, reusableCF, ct)
				}

			case reflect.Map:

				var pv string
				reusableCF := &cField{}

				for _, key := range current.MapKeys() {

					pv = fmt.Sprintf("%v", key.Interface())

					v.misc = append(v.misc[0:0], cf.name...)
					v.misc = append(v.misc, '[')
					v.misc = append(v.misc, pv...)
					v.misc = append(v.misc, ']')

					reusableCF.name = string(v.misc)

					if cf.namesEqual {
						reusableCF.altName = reusableCF.name
					} else {
						v.misc = append(v.misc[0:0], cf.altName...)
						v.misc = append(v.misc, '[')
						v.misc = append(v.misc, pv...)
						v.misc = append(v.misc, ']')

						reusableCF.altName = string(v.misc)
					}

					if ct != nil && ct.typeof == typeKeys && ct.keys != nil {
						v.traverseField(ctx, parent, key, ns, structNs, reusableCF, ct.keys)
						
						if ct.next != nil {
							v.traverseField(ctx, parent, current.MapIndex(key), ns, structNs, reusableCF, ct.next)
						}
					} else {
						v.traverseField(ctx, parent, current.MapIndex(key), ns, structNs, reusableCF, ct)
					}
				}

			default:
				
				
				panic("dive error! can't dive on a non slice or map")
			}

			return

		case typeOr:

			v.misc = v.misc[0:0]

			for {

				
				v.slflParent = parent
				v.flField = current
				v.cf = cf
				v.ct = ct

				if ct.fn(ctx, v) {
					if ct.isBlockEnd {
						ct = ct.next
						continue OUTER
					}

					
					for {

						ct = ct.next

						if ct == nil {
							continue OUTER
						}

						if ct.typeof != typeOr {
							continue OUTER
						}

						if ct.isBlockEnd {
							ct = ct.next
							continue OUTER
						}
					}
				}

				v.misc = append(v.misc, '|')
				v.misc = append(v.misc, ct.tag...)

				if ct.hasParam {
					v.misc = append(v.misc, '=')
					v.misc = append(v.misc, ct.param...)
				}

				if ct.isBlockEnd || ct.next == nil {
					
					v.str1 = string(append(ns, cf.altName...))

					if v.v.hasTagNameFunc {
						v.str2 = string(append(structNs, cf.name...))
					} else {
						v.str2 = v.str1
					}

					if ct.hasAlias {

						v.errs = append(v.errs,
							&fieldError{
								v:              v.v,
								tag:            ct.aliasTag,
								actualTag:      ct.actualAliasTag,
								ns:             v.str1,
								structNs:       v.str2,
								fieldLen:       uint8(len(cf.altName)),
								structfieldLen: uint8(len(cf.name)),
								value:          getValue(current),
								param:          ct.param,
								kind:           kind,
								typ:            typ,
							},
						)

					} else {

						tVal := string(v.misc)[1:]

						v.errs = append(v.errs,
							&fieldError{
								v:              v.v,
								tag:            tVal,
								actualTag:      tVal,
								ns:             v.str1,
								structNs:       v.str2,
								fieldLen:       uint8(len(cf.altName)),
								structfieldLen: uint8(len(cf.name)),
								value:          getValue(current),
								param:          ct.param,
								kind:           kind,
								typ:            typ,
							},
						)
					}

					return
				}

				ct = ct.next
			}

		default:

			
			v.slflParent = parent
			v.flField = current
			v.cf = cf
			v.ct = ct

			if !ct.fn(ctx, v) {
				v.str1 = string(append(ns, cf.altName...))

				if v.v.hasTagNameFunc {
					v.str2 = string(append(structNs, cf.name...))
				} else {
					v.str2 = v.str1
				}

				v.errs = append(v.errs,
					&fieldError{
						v:              v.v,
						tag:            ct.aliasTag,
						actualTag:      ct.tag,
						ns:             v.str1,
						structNs:       v.str2,
						fieldLen:       uint8(len(cf.altName)),
						structfieldLen: uint8(len(cf.name)),
						value:          getValue(current),
						param:          ct.param,
						kind:           kind,
						typ:            typ,
					},
				)

				return
			}
			ct = ct.next
		}
	}

}

func getValue(val reflect.Value) interface{} {
	if val.CanInterface() {
		return val.Interface()
	}

	if val.CanAddr() {
		return reflect.NewAt(val.Type(), unsafe.Pointer(val.UnsafeAddr())).Elem().Interface()
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint()
	case reflect.Complex64, reflect.Complex128:
		return val.Complex()
	case reflect.Float32, reflect.Float64:
		return val.Float()
	default:
		return val.String()
	}
}
