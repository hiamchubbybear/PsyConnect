package toml

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/pelletier/go-toml/v2/internal/danger"
	"github.com/pelletier/go-toml/v2/internal/tracker"
	"github.com/pelletier/go-toml/v2/unstable"
)




func Unmarshal(data []byte, v interface{}) error {
	d := decoder{}
	d.p.Reset(data)
	return d.FromParser(v)
}


type Decoder struct {
	
	r io.Reader

	
	strict bool

	
	unmarshalerInterface bool
}


func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}








func (d *Decoder) DisallowUnknownFields() *Decoder {
	d.strict = true
	return d
}














func (d *Decoder) EnableUnmarshalerInterface() *Decoder {
	d.unmarshalerInterface = true
	return d
}










































func (d *Decoder) Decode(v interface{}) error {
	b, err := io.ReadAll(d.r)
	if err != nil {
		return fmt.Errorf("toml: %w", err)
	}

	dec := decoder{
		strict: strict{
			Enabled: d.strict,
		},
		unmarshalerInterface: d.unmarshalerInterface,
	}
	dec.p.Reset(b)

	return dec.FromParser(v)
}

type decoder struct {
	
	p unstable.Parser

	
	
	
	stashedExpr bool

	
	
	
	skipUntilTable bool

	
	
	clearArrayTable bool

	
	
	
	
	arrayIndexes map[reflect.Value]int

	
	seen tracker.SeenTracker

	
	strict strict

	
	unmarshalerInterface bool

	
	errorContext *errorContext
}

type errorContext struct {
	Struct reflect.Type
	Field  []int
}

func (d *decoder) typeMismatchError(toml string, target reflect.Type) error {
	return fmt.Errorf("toml: %s", d.typeMismatchString(toml, target))
}

func (d *decoder) typeMismatchString(toml string, target reflect.Type) string {
	if d.errorContext != nil && d.errorContext.Struct != nil {
		ctx := d.errorContext
		f := ctx.Struct.FieldByIndex(ctx.Field)
		return fmt.Sprintf("cannot decode TOML %s into struct field %s.%s of type %s", toml, ctx.Struct, f.Name, f.Type)
	}
	return fmt.Sprintf("cannot decode TOML %s into a Go value of type %s", toml, target)
}

func (d *decoder) expr() *unstable.Node {
	return d.p.Expression()
}

func (d *decoder) nextExpr() bool {
	if d.stashedExpr {
		d.stashedExpr = false
		return true
	}
	return d.p.NextExpression()
}

func (d *decoder) stashExpr() {
	d.stashedExpr = true
}

func (d *decoder) arrayIndex(shouldAppend bool, v reflect.Value) int {
	if d.arrayIndexes == nil {
		d.arrayIndexes = make(map[reflect.Value]int, 1)
	}

	idx, ok := d.arrayIndexes[v]

	if !ok {
		d.arrayIndexes[v] = 0
	} else if shouldAppend {
		idx++
		d.arrayIndexes[v] = idx
	}

	return idx
}

func (d *decoder) FromParser(v interface{}) error {
	r := reflect.ValueOf(v)
	if r.Kind() != reflect.Ptr {
		return fmt.Errorf("toml: decoding can only be performed into a pointer, not %s", r.Kind())
	}

	if r.IsNil() {
		return fmt.Errorf("toml: decoding pointer target cannot be nil")
	}

	r = r.Elem()
	if r.Kind() == reflect.Interface && r.IsNil() {
		newMap := map[string]interface{}{}
		r.Set(reflect.ValueOf(newMap))
	}

	err := d.fromParser(r)
	if err == nil {
		return d.strict.Error(d.p.Data())
	}

	var e *unstable.ParserError
	if errors.As(err, &e) {
		return wrapDecodeError(d.p.Data(), e)
	}

	return err
}

func (d *decoder) fromParser(root reflect.Value) error {
	for d.nextExpr() {
		err := d.handleRootExpression(d.expr(), root)
		if err != nil {
			return err
		}
	}

	return d.p.Error()
}



func (d *decoder) handleRootExpression(expr *unstable.Node, v reflect.Value) error {
	var x reflect.Value
	var err error
	var first bool 

	if !(d.skipUntilTable && expr.Kind == unstable.KeyValue) {
		first, err = d.seen.CheckExpression(expr)
		if err != nil {
			return err
		}
	}

	switch expr.Kind {
	case unstable.KeyValue:
		if d.skipUntilTable {
			return nil
		}
		x, err = d.handleKeyValue(expr, v)
	case unstable.Table:
		d.skipUntilTable = false
		d.strict.EnterTable(expr)
		x, err = d.handleTable(expr.Key(), v)
	case unstable.ArrayTable:
		d.skipUntilTable = false
		d.strict.EnterArrayTable(expr)
		d.clearArrayTable = first
		x, err = d.handleArrayTable(expr.Key(), v)
	default:
		panic(fmt.Errorf("parser should not permit expression of kind %s at document root", expr.Kind))
	}

	if d.skipUntilTable {
		if expr.Kind == unstable.Table || expr.Kind == unstable.ArrayTable {
			d.strict.MissingTable(expr)
		}
	} else if err == nil && x.IsValid() {
		v.Set(x)
	}

	return err
}

func (d *decoder) handleArrayTable(key unstable.Iterator, v reflect.Value) (reflect.Value, error) {
	if key.Next() {
		return d.handleArrayTablePart(key, v)
	}
	return d.handleKeyValues(v)
}

func (d *decoder) handleArrayTableCollectionLast(key unstable.Iterator, v reflect.Value) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Interface:
		elem := v.Elem()
		if !elem.IsValid() {
			elem = reflect.New(sliceInterfaceType).Elem()
			elem.Set(reflect.MakeSlice(sliceInterfaceType, 0, 16))
		} else if elem.Kind() == reflect.Slice {
			if elem.Type() != sliceInterfaceType {
				elem = reflect.New(sliceInterfaceType).Elem()
				elem.Set(reflect.MakeSlice(sliceInterfaceType, 0, 16))
			} else if !elem.CanSet() {
				nelem := reflect.New(sliceInterfaceType).Elem()
				nelem.Set(reflect.MakeSlice(sliceInterfaceType, elem.Len(), elem.Cap()))
				reflect.Copy(nelem, elem)
				elem = nelem
			}
			if d.clearArrayTable && elem.Len() > 0 {
				elem.SetLen(0)
				d.clearArrayTable = false
			}
		}
		return d.handleArrayTableCollectionLast(key, elem)
	case reflect.Ptr:
		elem := v.Elem()
		if !elem.IsValid() {
			ptr := reflect.New(v.Type().Elem())
			v.Set(ptr)
			elem = ptr.Elem()
		}

		elem, err := d.handleArrayTableCollectionLast(key, elem)
		if err != nil {
			return reflect.Value{}, err
		}
		v.Elem().Set(elem)

		return v, nil
	case reflect.Slice:
		if d.clearArrayTable && v.Len() > 0 {
			v.SetLen(0)
			d.clearArrayTable = false
		}
		elemType := v.Type().Elem()
		var elem reflect.Value
		if elemType.Kind() == reflect.Interface {
			elem = makeMapStringInterface()
		} else {
			elem = reflect.New(elemType).Elem()
		}
		elem2, err := d.handleArrayTable(key, elem)
		if err != nil {
			return reflect.Value{}, err
		}
		if elem2.IsValid() {
			elem = elem2
		}
		return reflect.Append(v, elem), nil
	case reflect.Array:
		idx := d.arrayIndex(true, v)
		if idx >= v.Len() {
			return v, fmt.Errorf("%s at position %d", d.typeMismatchError("array table", v.Type()), idx)
		}
		elem := v.Index(idx)
		_, err := d.handleArrayTable(key, elem)
		return v, err
	default:
		return reflect.Value{}, d.typeMismatchError("array table", v.Type())
	}
}





func (d *decoder) handleArrayTableCollection(key unstable.Iterator, v reflect.Value) (reflect.Value, error) {
	if key.IsLast() {
		return d.handleArrayTableCollectionLast(key, v)
	}

	switch v.Kind() {
	case reflect.Ptr:
		elem := v.Elem()
		if !elem.IsValid() {
			ptr := reflect.New(v.Type().Elem())
			v.Set(ptr)
			elem = ptr.Elem()
		}

		elem, err := d.handleArrayTableCollection(key, elem)
		if err != nil {
			return reflect.Value{}, err
		}
		if elem.IsValid() {
			v.Elem().Set(elem)
		}

		return v, nil
	case reflect.Slice:
		elem := v.Index(v.Len() - 1)
		x, err := d.handleArrayTable(key, elem)
		if err != nil || d.skipUntilTable {
			return reflect.Value{}, err
		}
		if x.IsValid() {
			elem.Set(x)
		}

		return v, err
	case reflect.Array:
		idx := d.arrayIndex(false, v)
		if idx >= v.Len() {
			return v, fmt.Errorf("%s at position %d", d.typeMismatchError("array table", v.Type()), idx)
		}
		elem := v.Index(idx)
		_, err := d.handleArrayTable(key, elem)
		return v, err
	}

	return d.handleArrayTable(key, v)
}

func (d *decoder) handleKeyPart(key unstable.Iterator, v reflect.Value, nextFn handlerFn, makeFn valueMakerFn) (reflect.Value, error) {
	var rv reflect.Value

	
	
	switch v.Kind() {
	case reflect.Ptr:
		elem := v.Elem()
		if !elem.IsValid() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		elem = v.Elem()
		return d.handleKeyPart(key, elem, nextFn, makeFn)
	case reflect.Map:
		vt := v.Type()

		
		mk, err := d.keyFromData(vt.Key(), key.Node().Data)
		if err != nil {
			return reflect.Value{}, err
		}

		
		if v.IsNil() {
			vt := v.Type()
			v = reflect.MakeMap(vt)
			rv = v
		}

		mv := v.MapIndex(mk)
		set := false
		if !mv.IsValid() {
			
			
			
			

			t := vt.Elem()
			if t.Kind() == reflect.Interface {
				mv = makeFn()
			} else {
				mv = reflect.New(t).Elem()
			}
			set = true
		} else if mv.Kind() == reflect.Interface {
			mv = mv.Elem()
			if !mv.IsValid() {
				mv = makeFn()
			}
			set = true
		} else if !mv.CanAddr() {
			vt := v.Type()
			t := vt.Elem()
			oldmv := mv
			mv = reflect.New(t).Elem()
			mv.Set(oldmv)
			set = true
		}

		x, err := nextFn(key, mv)
		if err != nil {
			return reflect.Value{}, err
		}

		if x.IsValid() {
			mv = x
			set = true
		}

		if set {
			v.SetMapIndex(mk, mv)
		}
	case reflect.Struct:
		path, found := structFieldPath(v, string(key.Node().Data))
		if !found {
			d.skipUntilTable = true
			return reflect.Value{}, nil
		}

		if d.errorContext == nil {
			d.errorContext = new(errorContext)
		}
		t := v.Type()
		d.errorContext.Struct = t
		d.errorContext.Field = path

		f := fieldByIndex(v, path)
		x, err := nextFn(key, f)
		if err != nil || d.skipUntilTable {
			return reflect.Value{}, err
		}
		if x.IsValid() {
			f.Set(x)
		}
		d.errorContext.Field = nil
		d.errorContext.Struct = nil
	case reflect.Interface:
		if v.Elem().IsValid() {
			v = v.Elem()
		} else {
			v = makeMapStringInterface()
		}

		x, err := d.handleKeyPart(key, v, nextFn, makeFn)
		if err != nil {
			return reflect.Value{}, err
		}
		if x.IsValid() {
			v = x
		}
		rv = v
	default:
		panic(fmt.Errorf("unhandled part: %s", v.Kind()))
	}

	return rv, nil
}




func (d *decoder) handleArrayTablePart(key unstable.Iterator, v reflect.Value) (reflect.Value, error) {
	var makeFn valueMakerFn
	if key.IsLast() {
		makeFn = makeSliceInterface
	} else {
		makeFn = makeMapStringInterface
	}
	return d.handleKeyPart(key, v, d.handleArrayTableCollection, makeFn)
}



func (d *decoder) handleTable(key unstable.Iterator, v reflect.Value) (reflect.Value, error) {
	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			return reflect.Value{}, unstable.NewParserError(key.Node().Data, "cannot store a table in a slice")
		}
		elem := v.Index(v.Len() - 1)
		x, err := d.handleTable(key, elem)
		if err != nil {
			return reflect.Value{}, err
		}
		if x.IsValid() {
			elem.Set(x)
		}
		return reflect.Value{}, nil
	}
	if key.Next() {
		
		return d.handleTablePart(key, v)
	}
	
	
	return d.handleKeyValues(v)
}



func (d *decoder) handleKeyValues(v reflect.Value) (reflect.Value, error) {
	var rv reflect.Value
	for d.nextExpr() {
		expr := d.expr()
		if expr.Kind != unstable.KeyValue {
			
			
			
			
			d.stashExpr()
			break
		}

		_, err := d.seen.CheckExpression(expr)
		if err != nil {
			return reflect.Value{}, err
		}

		x, err := d.handleKeyValue(expr, v)
		if err != nil {
			return reflect.Value{}, err
		}
		if x.IsValid() {
			v = x
			rv = x
		}
	}
	return rv, nil
}

type (
	handlerFn    func(key unstable.Iterator, v reflect.Value) (reflect.Value, error)
	valueMakerFn func() reflect.Value
)

func makeMapStringInterface() reflect.Value {
	return reflect.MakeMap(mapStringInterfaceType)
}

func makeSliceInterface() reflect.Value {
	return reflect.MakeSlice(sliceInterfaceType, 0, 16)
}

func (d *decoder) handleTablePart(key unstable.Iterator, v reflect.Value) (reflect.Value, error) {
	return d.handleKeyPart(key, v, d.handleTable, makeMapStringInterface)
}

func (d *decoder) tryTextUnmarshaler(node *unstable.Node, v reflect.Value) (bool, error) {
	
	
	if v.Type() == timeType {
		return false, nil
	}

	if v.CanAddr() && v.Addr().Type().Implements(textUnmarshalerType) {
		err := v.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText(node.Data)
		if err != nil {
			return false, unstable.NewParserError(d.p.Raw(node.Raw), "%w", err)
		}

		return true, nil
	}

	return false, nil
}

func (d *decoder) handleValue(value *unstable.Node, v reflect.Value) error {
	for v.Kind() == reflect.Ptr {
		v = initAndDereferencePointer(v)
	}

	if d.unmarshalerInterface {
		if v.CanAddr() && v.Addr().CanInterface() {
			if outi, ok := v.Addr().Interface().(unstable.Unmarshaler); ok {
				return outi.UnmarshalTOML(value)
			}
		}
	}

	ok, err := d.tryTextUnmarshaler(value, v)
	if ok || err != nil {
		return err
	}

	switch value.Kind {
	case unstable.String:
		return d.unmarshalString(value, v)
	case unstable.Integer:
		return d.unmarshalInteger(value, v)
	case unstable.Float:
		return d.unmarshalFloat(value, v)
	case unstable.Bool:
		return d.unmarshalBool(value, v)
	case unstable.DateTime:
		return d.unmarshalDateTime(value, v)
	case unstable.LocalDate:
		return d.unmarshalLocalDate(value, v)
	case unstable.LocalTime:
		return d.unmarshalLocalTime(value, v)
	case unstable.LocalDateTime:
		return d.unmarshalLocalDateTime(value, v)
	case unstable.InlineTable:
		return d.unmarshalInlineTable(value, v)
	case unstable.Array:
		return d.unmarshalArray(value, v)
	default:
		panic(fmt.Errorf("handleValue not implemented for %s", value.Kind))
	}
}

func (d *decoder) unmarshalArray(array *unstable.Node, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Slice:
		if v.IsNil() {
			v.Set(reflect.MakeSlice(v.Type(), 0, 16))
		} else {
			v.SetLen(0)
		}
	case reflect.Array:
		
	case reflect.Interface:
		elem := v.Elem()
		if !elem.IsValid() {
			elem = reflect.New(sliceInterfaceType).Elem()
			elem.Set(reflect.MakeSlice(sliceInterfaceType, 0, 16))
		} else if elem.Kind() == reflect.Slice {
			if elem.Type() != sliceInterfaceType {
				elem = reflect.New(sliceInterfaceType).Elem()
				elem.Set(reflect.MakeSlice(sliceInterfaceType, 0, 16))
			} else if !elem.CanSet() {
				nelem := reflect.New(sliceInterfaceType).Elem()
				nelem.Set(reflect.MakeSlice(sliceInterfaceType, elem.Len(), elem.Cap()))
				reflect.Copy(nelem, elem)
				elem = nelem
			}
		}
		err := d.unmarshalArray(array, elem)
		if err != nil {
			return err
		}
		v.Set(elem)
		return nil
	default:
		
		
		return d.typeMismatchError("array", v.Type())
	}

	elemType := v.Type().Elem()

	it := array.Children()
	idx := 0
	for it.Next() {
		n := it.Node()

		
		if v.Kind() == reflect.Slice {
			elem := reflect.New(elemType).Elem()

			err := d.handleValue(n, elem)
			if err != nil {
				return err
			}

			v.Set(reflect.Append(v, elem))
		} else { 
			if idx >= v.Len() {
				return nil
			}
			elem := v.Index(idx)
			err := d.handleValue(n, elem)
			if err != nil {
				return err
			}
			idx++
		}
	}

	return nil
}

func (d *decoder) unmarshalInlineTable(itable *unstable.Node, v reflect.Value) error {
	
	switch v.Kind() {
	case reflect.Map:
		if v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
		}
	case reflect.Struct:
	
	case reflect.Interface:
		elem := v.Elem()
		if !elem.IsValid() {
			elem = makeMapStringInterface()
			v.Set(elem)
		}
		return d.unmarshalInlineTable(itable, elem)
	default:
		return unstable.NewParserError(d.p.Raw(itable.Raw), "cannot store inline table in Go type %s", v.Kind())
	}

	it := itable.Children()
	for it.Next() {
		n := it.Node()

		x, err := d.handleKeyValue(n, v)
		if err != nil {
			return err
		}
		if x.IsValid() {
			v = x
		}
	}

	return nil
}

func (d *decoder) unmarshalDateTime(value *unstable.Node, v reflect.Value) error {
	dt, err := parseDateTime(value.Data)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(dt))
	return nil
}

func (d *decoder) unmarshalLocalDate(value *unstable.Node, v reflect.Value) error {
	ld, err := parseLocalDate(value.Data)
	if err != nil {
		return err
	}

	if v.Type() == timeType {
		cast := ld.AsTime(time.Local)
		v.Set(reflect.ValueOf(cast))
		return nil
	}

	v.Set(reflect.ValueOf(ld))

	return nil
}

func (d *decoder) unmarshalLocalTime(value *unstable.Node, v reflect.Value) error {
	lt, rest, err := parseLocalTime(value.Data)
	if err != nil {
		return err
	}

	if len(rest) > 0 {
		return unstable.NewParserError(rest, "extra characters at the end of a local time")
	}

	v.Set(reflect.ValueOf(lt))
	return nil
}

func (d *decoder) unmarshalLocalDateTime(value *unstable.Node, v reflect.Value) error {
	ldt, rest, err := parseLocalDateTime(value.Data)
	if err != nil {
		return err
	}

	if len(rest) > 0 {
		return unstable.NewParserError(rest, "extra characters at the end of a local date time")
	}

	if v.Type() == timeType {
		cast := ldt.AsTime(time.Local)

		v.Set(reflect.ValueOf(cast))
		return nil
	}

	v.Set(reflect.ValueOf(ldt))

	return nil
}

func (d *decoder) unmarshalBool(value *unstable.Node, v reflect.Value) error {
	b := value.Data[0] == 't'

	switch v.Kind() {
	case reflect.Bool:
		v.SetBool(b)
	case reflect.Interface:
		v.Set(reflect.ValueOf(b))
	default:
		return unstable.NewParserError(value.Data, "cannot assign boolean to a %t", b)
	}

	return nil
}

func (d *decoder) unmarshalFloat(value *unstable.Node, v reflect.Value) error {
	f, err := parseFloat(value.Data)
	if err != nil {
		return err
	}

	switch v.Kind() {
	case reflect.Float64:
		v.SetFloat(f)
	case reflect.Float32:
		if f > math.MaxFloat32 {
			return unstable.NewParserError(value.Data, "number %f does not fit in a float32", f)
		}
		v.SetFloat(f)
	case reflect.Interface:
		v.Set(reflect.ValueOf(f))
	default:
		return unstable.NewParserError(value.Data, "float cannot be assigned to %s", v.Kind())
	}

	return nil
}

const (
	maxInt = int64(^uint(0) >> 1)
	minInt = -maxInt - 1
)







var maxUint int64 = math.MaxInt64

func init() {
	m := uint64(^uint(0))
	if m < uint64(maxUint) {
		maxUint = int64(m)
	}
}

func (d *decoder) unmarshalInteger(value *unstable.Node, v reflect.Value) error {
	kind := v.Kind()
	if kind == reflect.Float32 || kind == reflect.Float64 {
		return d.unmarshalFloat(value, v)
	}

	i, err := parseInteger(value.Data)
	if err != nil {
		return err
	}

	var r reflect.Value

	switch kind {
	case reflect.Int64:
		v.SetInt(i)
		return nil
	case reflect.Int32:
		if i < math.MinInt32 || i > math.MaxInt32 {
			return fmt.Errorf("toml: number %d does not fit in an int32", i)
		}

		r = reflect.ValueOf(int32(i))
	case reflect.Int16:
		if i < math.MinInt16 || i > math.MaxInt16 {
			return fmt.Errorf("toml: number %d does not fit in an int16", i)
		}

		r = reflect.ValueOf(int16(i))
	case reflect.Int8:
		if i < math.MinInt8 || i > math.MaxInt8 {
			return fmt.Errorf("toml: number %d does not fit in an int8", i)
		}

		r = reflect.ValueOf(int8(i))
	case reflect.Int:
		if i < minInt || i > maxInt {
			return fmt.Errorf("toml: number %d does not fit in an int", i)
		}

		r = reflect.ValueOf(int(i))
	case reflect.Uint64:
		if i < 0 {
			return fmt.Errorf("toml: negative number %d does not fit in an uint64", i)
		}

		r = reflect.ValueOf(uint64(i))
	case reflect.Uint32:
		if i < 0 || i > math.MaxUint32 {
			return fmt.Errorf("toml: negative number %d does not fit in an uint32", i)
		}

		r = reflect.ValueOf(uint32(i))
	case reflect.Uint16:
		if i < 0 || i > math.MaxUint16 {
			return fmt.Errorf("toml: negative number %d does not fit in an uint16", i)
		}

		r = reflect.ValueOf(uint16(i))
	case reflect.Uint8:
		if i < 0 || i > math.MaxUint8 {
			return fmt.Errorf("toml: negative number %d does not fit in an uint8", i)
		}

		r = reflect.ValueOf(uint8(i))
	case reflect.Uint:
		if i < 0 || i > maxUint {
			return fmt.Errorf("toml: negative number %d does not fit in an uint", i)
		}

		r = reflect.ValueOf(uint(i))
	case reflect.Interface:
		r = reflect.ValueOf(i)
	default:
		return unstable.NewParserError(d.p.Raw(value.Raw), d.typeMismatchString("integer", v.Type()))
	}

	if !r.Type().AssignableTo(v.Type()) {
		r = r.Convert(v.Type())
	}

	v.Set(r)

	return nil
}

func (d *decoder) unmarshalString(value *unstable.Node, v reflect.Value) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(string(value.Data))
	case reflect.Interface:
		v.Set(reflect.ValueOf(string(value.Data)))
	default:
		return unstable.NewParserError(d.p.Raw(value.Raw), d.typeMismatchString("string", v.Type()))
	}

	return nil
}

func (d *decoder) handleKeyValue(expr *unstable.Node, v reflect.Value) (reflect.Value, error) {
	d.strict.EnterKeyValue(expr)

	v, err := d.handleKeyValueInner(expr.Key(), expr.Value(), v)
	if d.skipUntilTable {
		d.strict.MissingField(expr)
		d.skipUntilTable = false
	}

	d.strict.ExitKeyValue(expr)

	return v, err
}

func (d *decoder) handleKeyValueInner(key unstable.Iterator, value *unstable.Node, v reflect.Value) (reflect.Value, error) {
	if key.Next() {
		
		return d.handleKeyValuePart(key, value, v)
	}
	
	
	return reflect.Value{}, d.handleValue(value, v)
}

func (d *decoder) keyFromData(keyType reflect.Type, data []byte) (reflect.Value, error) {
	switch {
	case stringType.AssignableTo(keyType):
		return reflect.ValueOf(string(data)), nil

	case stringType.ConvertibleTo(keyType):
		return reflect.ValueOf(string(data)).Convert(keyType), nil

	case keyType.Implements(textUnmarshalerType):
		mk := reflect.New(keyType.Elem())
		if err := mk.Interface().(encoding.TextUnmarshaler).UnmarshalText(data); err != nil {
			return reflect.Value{}, fmt.Errorf("toml: error unmarshalling key type %s from text: %w", stringType, err)
		}
		return mk, nil

	case reflect.PointerTo(keyType).Implements(textUnmarshalerType):
		mk := reflect.New(keyType)
		if err := mk.Interface().(encoding.TextUnmarshaler).UnmarshalText(data); err != nil {
			return reflect.Value{}, fmt.Errorf("toml: error unmarshalling key type %s from text: %w", stringType, err)
		}
		return mk.Elem(), nil

	case keyType.Kind() == reflect.Int || keyType.Kind() == reflect.Int8 || keyType.Kind() == reflect.Int16 || keyType.Kind() == reflect.Int32 || keyType.Kind() == reflect.Int64:
		key, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("toml: error parsing key of type %s from integer: %w", stringType, err)
		}
		return reflect.ValueOf(key).Convert(keyType), nil
	case keyType.Kind() == reflect.Uint || keyType.Kind() == reflect.Uint8 || keyType.Kind() == reflect.Uint16 || keyType.Kind() == reflect.Uint32 || keyType.Kind() == reflect.Uint64:
		key, err := strconv.ParseUint(string(data), 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("toml: error parsing key of type %s from unsigned integer: %w", stringType, err)
		}
		return reflect.ValueOf(key).Convert(keyType), nil

	case keyType.Kind() == reflect.Float32:
		key, err := strconv.ParseFloat(string(data), 32)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("toml: error parsing key of type %s from float: %w", stringType, err)
		}
		return reflect.ValueOf(float32(key)), nil

	case keyType.Kind() == reflect.Float64:
		key, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("toml: error parsing key of type %s from float: %w", stringType, err)
		}
		return reflect.ValueOf(float64(key)), nil
	}
	return reflect.Value{}, fmt.Errorf("toml: cannot convert map key of type %s to expected type %s", stringType, keyType)
}

func (d *decoder) handleKeyValuePart(key unstable.Iterator, value *unstable.Node, v reflect.Value) (reflect.Value, error) {
	
	var rv reflect.Value

	
	
	switch v.Kind() {
	case reflect.Map:
		vt := v.Type()

		mk, err := d.keyFromData(vt.Key(), key.Node().Data)
		if err != nil {
			return reflect.Value{}, err
		}

		
		if v.IsNil() {
			v = reflect.MakeMap(vt)
			rv = v
		}

		mv := v.MapIndex(mk)
		set := false
		if !mv.IsValid() || key.IsLast() {
			set = true
			mv = reflect.New(v.Type().Elem()).Elem()
		}

		nv, err := d.handleKeyValueInner(key, value, mv)
		if err != nil {
			return reflect.Value{}, err
		}
		if nv.IsValid() {
			mv = nv
			set = true
		}

		if set {
			v.SetMapIndex(mk, mv)
		}
	case reflect.Struct:
		path, found := structFieldPath(v, string(key.Node().Data))
		if !found {
			d.skipUntilTable = true
			break
		}

		if d.errorContext == nil {
			d.errorContext = new(errorContext)
		}
		t := v.Type()
		d.errorContext.Struct = t
		d.errorContext.Field = path

		f := fieldByIndex(v, path)

		if !f.CanAddr() {
			
			
			nvp := reflect.New(v.Type())
			nvp.Elem().Set(v)
			v = nvp.Elem()
			_, err := d.handleKeyValuePart(key, value, v)
			if err != nil {
				return reflect.Value{}, err
			}
			return nvp.Elem(), nil
		}
		x, err := d.handleKeyValueInner(key, value, f)
		if err != nil {
			return reflect.Value{}, err
		}

		if x.IsValid() {
			f.Set(x)
		}
		d.errorContext.Struct = nil
		d.errorContext.Field = nil
	case reflect.Interface:
		v = v.Elem()

		
		
		
		
		if !v.IsValid() || v.Type() != mapStringInterfaceType {
			v = makeMapStringInterface()
		}

		x, err := d.handleKeyValuePart(key, value, v)
		if err != nil {
			return reflect.Value{}, err
		}
		if x.IsValid() {
			v = x
		}
		rv = v
	case reflect.Ptr:
		elem := v.Elem()
		if !elem.IsValid() {
			ptr := reflect.New(v.Type().Elem())
			v.Set(ptr)
			rv = v
			elem = ptr.Elem()
		}

		elem2, err := d.handleKeyValuePart(key, value, elem)
		if err != nil {
			return reflect.Value{}, err
		}
		if elem2.IsValid() {
			elem = elem2
		}
		v.Elem().Set(elem)
	default:
		return reflect.Value{}, fmt.Errorf("unhandled kv part: %s", v.Kind())
	}

	return rv, nil
}

func initAndDereferencePointer(v reflect.Value) reflect.Value {
	var elem reflect.Value
	if v.IsNil() {
		ptr := reflect.New(v.Type().Elem())
		v.Set(ptr)
	}
	elem = v.Elem()
	return elem
}


func fieldByIndex(v reflect.Value, path []int) reflect.Value {
	for _, x := range path {
		v = v.Field(x)

		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		}
	}
	return v
}

type fieldPathsMap = map[string][]int

var globalFieldPathsCache atomic.Value 

func structFieldPath(v reflect.Value, name string) ([]int, bool) {
	t := v.Type()

	cache, _ := globalFieldPathsCache.Load().(map[danger.TypeID]fieldPathsMap)
	fieldPaths, ok := cache[danger.MakeTypeID(t)]

	if !ok {
		fieldPaths = map[string][]int{}

		forEachField(t, nil, func(name string, path []int) {
			fieldPaths[name] = path
			
			fieldPaths[strings.ToLower(name)] = path
		})

		newCache := make(map[danger.TypeID]fieldPathsMap, len(cache)+1)
		newCache[danger.MakeTypeID(t)] = fieldPaths
		for k, v := range cache {
			newCache[k] = v
		}
		globalFieldPathsCache.Store(newCache)
	}

	path, ok := fieldPaths[name]
	if !ok {
		path, ok = fieldPaths[strings.ToLower(name)]
	}
	return path, ok
}

func forEachField(t reflect.Type, path []int, do func(name string, path []int)) {
	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)

		if !f.Anonymous && f.PkgPath != "" {
			
			continue
		}

		fieldPath := append(path, i)
		fieldPath = fieldPath[:len(fieldPath):len(fieldPath)]

		name := f.Tag.Get("toml")
		if name == "-" {
			continue
		}

		if i := strings.IndexByte(name, ','); i >= 0 {
			name = name[:i]
		}

		if f.Anonymous && name == "" {
			t2 := f.Type
			if t2.Kind() == reflect.Ptr {
				t2 = t2.Elem()
			}

			if t2.Kind() == reflect.Struct {
				forEachField(t2, fieldPath, do)
			}
			continue
		}

		if name == "" {
			name = f.Name
		}

		do(name, fieldPath)
	}
}
