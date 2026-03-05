package validator

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	ut "github.com/go-playground/universal-translator"
)

const (
	defaultTagName        = "validate"
	utf8HexComma          = "0x2C"
	utf8Pipe              = "0x7C"
	tagSeparator          = ","
	orSeparator           = "|"
	tagKeySeparator       = "="
	structOnlyTag         = "structonly"
	noStructLevelTag      = "nostructlevel"
	omitzero              = "omitzero"
	omitempty             = "omitempty"
	omitnil               = "omitnil"
	isdefault             = "isdefault"
	requiredWithoutAllTag = "required_without_all"
	requiredWithoutTag    = "required_without"
	requiredWithTag       = "required_with"
	requiredWithAllTag    = "required_with_all"
	requiredIfTag         = "required_if"
	requiredUnlessTag     = "required_unless"
	skipUnlessTag         = "skip_unless"
	excludedWithoutAllTag = "excluded_without_all"
	excludedWithoutTag    = "excluded_without"
	excludedWithTag       = "excluded_with"
	excludedWithAllTag    = "excluded_with_all"
	excludedIfTag         = "excluded_if"
	excludedUnlessTag     = "excluded_unless"
	skipValidationTag     = "-"
	diveTag               = "dive"
	keysTag               = "keys"
	endKeysTag            = "endkeys"
	requiredTag           = "required"
	namespaceSeparator    = "."
	leftBracket           = "["
	rightBracket          = "]"
	restrictedTagChars    = ".[],|=+()`~!@#$%^&*\\\"/?<>{}"
	restrictedAliasErr    = "Alias '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
	restrictedTagErr      = "Tag '%s' either contains restricted characters or is the same as a restricted tag needed for normal operation"
)

var (
	timeDurationType = reflect.TypeOf(time.Duration(0))
	timeType         = reflect.TypeOf(time.Time{})

	byteSliceType = reflect.TypeOf([]byte{})

	defaultCField = &cField{namesEqual: true}
)





type FilterFunc func(ns []byte) bool




type CustomTypeFunc func(field reflect.Value) interface{}


type TagNameFunc func(field reflect.StructField) string

type internalValidationFuncWrapper struct {
	fn                 FuncCtx
	runValidationOnNil bool
}


type Validate struct {
	tagName                string
	pool                   *sync.Pool
	tagNameFunc            TagNameFunc
	structLevelFuncs       map[reflect.Type]StructLevelFuncCtx
	customFuncs            map[reflect.Type]CustomTypeFunc
	aliases                map[string]string
	validations            map[string]internalValidationFuncWrapper
	transTagFunc           map[ut.Translator]map[string]TranslationFunc 
	rules                  map[reflect.Type]map[string]string
	tagCache               *tagCache
	structCache            *structCache
	hasCustomFuncs         bool
	hasTagNameFunc         bool
	requiredStructEnabled  bool
	privateFieldValidation bool
}






func New(options ...Option) *Validate {

	tc := new(tagCache)
	tc.m.Store(make(map[string]*cTag))

	sc := new(structCache)
	sc.m.Store(make(map[reflect.Type]*cStruct))

	v := &Validate{
		tagName:     defaultTagName,
		aliases:     make(map[string]string, len(bakedInAliases)),
		validations: make(map[string]internalValidationFuncWrapper, len(bakedInValidators)),
		tagCache:    tc,
		structCache: sc,
	}

	
	for k, val := range bakedInAliases {
		v.RegisterAlias(k, val)
	}

	
	for k, val := range bakedInValidators {

		switch k {
		
		case requiredIfTag, requiredUnlessTag, requiredWithTag, requiredWithAllTag, requiredWithoutTag, requiredWithoutAllTag,
			excludedIfTag, excludedUnlessTag, excludedWithTag, excludedWithAllTag, excludedWithoutTag, excludedWithoutAllTag,
			skipUnlessTag:
			_ = v.registerValidation(k, wrapFunc(val), true, true)
		default:
			
			_ = v.registerValidation(k, wrapFunc(val), true, false)
		}
	}

	v.pool = &sync.Pool{
		New: func() interface{} {
			return &validate{
				v:        v,
				ns:       make([]byte, 0, 64),
				actualNs: make([]byte, 0, 64),
				misc:     make([]byte, 32),
			}
		},
	}

	for _, o := range options {
		o(v)
	}
	return v
}


func (v *Validate) SetTagName(name string) {
	v.tagName = name
}



func (v Validate) ValidateMapCtx(ctx context.Context, data map[string]interface{}, rules map[string]interface{}) map[string]interface{} {
	errs := make(map[string]interface{})
	for field, rule := range rules {
		if ruleObj, ok := rule.(map[string]interface{}); ok {
			if dataObj, ok := data[field].(map[string]interface{}); ok {
				err := v.ValidateMapCtx(ctx, dataObj, ruleObj)
				if len(err) > 0 {
					errs[field] = err
				}
			} else if dataObjs, ok := data[field].([]map[string]interface{}); ok {
				for _, obj := range dataObjs {
					err := v.ValidateMapCtx(ctx, obj, ruleObj)
					if len(err) > 0 {
						errs[field] = err
					}
				}
			} else {
				errs[field] = errors.New("The field: '" + field + "' is not a map to dive")
			}
		} else if ruleStr, ok := rule.(string); ok {
			err := v.VarCtx(ctx, data[field], ruleStr)
			if err != nil {
				errs[field] = err
			}
		}
	}
	return errs
}


func (v *Validate) ValidateMap(data map[string]interface{}, rules map[string]interface{}) map[string]interface{} {
	return v.ValidateMapCtx(context.Background(), data, rules)
}













func (v *Validate) RegisterTagNameFunc(fn TagNameFunc) {
	v.tagNameFunc = fn
	v.hasTagNameFunc = true
}






func (v *Validate) RegisterValidation(tag string, fn Func, callValidationEvenIfNull ...bool) error {
	return v.RegisterValidationCtx(tag, wrapFunc(fn), callValidationEvenIfNull...)
}



func (v *Validate) RegisterValidationCtx(tag string, fn FuncCtx, callValidationEvenIfNull ...bool) error {
	var nilCheckable bool
	if len(callValidationEvenIfNull) > 0 {
		nilCheckable = callValidationEvenIfNull[0]
	}
	return v.registerValidation(tag, fn, false, nilCheckable)
}

func (v *Validate) registerValidation(tag string, fn FuncCtx, bakedIn bool, nilCheckable bool) error {
	if len(tag) == 0 {
		return errors.New("function Key cannot be empty")
	}

	if fn == nil {
		return errors.New("function cannot be empty")
	}

	_, ok := restrictedTags[tag]
	if !bakedIn && (ok || strings.ContainsAny(tag, restrictedTagChars)) {
		panic(fmt.Sprintf(restrictedTagErr, tag))
	}
	v.validations[tag] = internalValidationFuncWrapper{fn: fn, runValidationOnNil: nilCheckable}
	return nil
}






func (v *Validate) RegisterAlias(alias, tags string) {

	_, ok := restrictedTags[alias]

	if ok || strings.ContainsAny(alias, restrictedTagChars) {
		panic(fmt.Sprintf(restrictedAliasErr, alias))
	}

	v.aliases[alias] = tags
}





func (v *Validate) RegisterStructValidation(fn StructLevelFunc, types ...interface{}) {
	v.RegisterStructValidationCtx(wrapStructLevelFunc(fn), types...)
}






func (v *Validate) RegisterStructValidationCtx(fn StructLevelFuncCtx, types ...interface{}) {

	if v.structLevelFuncs == nil {
		v.structLevelFuncs = make(map[reflect.Type]StructLevelFuncCtx)
	}

	for _, t := range types {
		tv := reflect.ValueOf(t)
		if tv.Kind() == reflect.Ptr {
			t = reflect.Indirect(tv).Interface()
		}

		v.structLevelFuncs[reflect.TypeOf(t)] = fn
	}
}





func (v *Validate) RegisterStructValidationMapRules(rules map[string]string, types ...interface{}) {
	if v.rules == nil {
		v.rules = make(map[reflect.Type]map[string]string)
	}

	deepCopyRules := make(map[string]string)
	for i, rule := range rules {
		deepCopyRules[i] = rule
	}

	for _, t := range types {
		typ := reflect.TypeOf(t)

		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		if typ.Kind() != reflect.Struct {
			continue
		}
		v.rules[typ] = deepCopyRules
	}
}




func (v *Validate) RegisterCustomTypeFunc(fn CustomTypeFunc, types ...interface{}) {

	if v.customFuncs == nil {
		v.customFuncs = make(map[reflect.Type]CustomTypeFunc)
	}

	for _, t := range types {
		v.customFuncs[reflect.TypeOf(t)] = fn
	}

	v.hasCustomFuncs = true
}


func (v *Validate) RegisterTranslation(tag string, trans ut.Translator, registerFn RegisterTranslationsFunc, translationFn TranslationFunc) (err error) {

	if v.transTagFunc == nil {
		v.transTagFunc = make(map[ut.Translator]map[string]TranslationFunc)
	}

	if err = registerFn(trans); err != nil {
		return
	}

	m, ok := v.transTagFunc[trans]
	if !ok {
		m = make(map[string]TranslationFunc)
		v.transTagFunc[trans] = m
	}

	m[tag] = translationFn

	return
}





func (v *Validate) Struct(s interface{}) error {
	return v.StructCtx(context.Background(), s)
}






func (v *Validate) StructCtx(ctx context.Context, s interface{}) (err error) {

	val := reflect.ValueOf(s)
	top := val

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct || val.Type().ConvertibleTo(timeType) {
		return &InvalidValidationError{Type: reflect.TypeOf(s)}
	}

	
	vd := v.pool.Get().(*validate)
	vd.top = top
	vd.isPartial = false
	

	vd.validateStruct(ctx, top, val, val.Type(), vd.ns[0:0], vd.actualNs[0:0], nil)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}

	v.pool.Put(vd)

	return
}






func (v *Validate) StructFiltered(s interface{}, fn FilterFunc) error {
	return v.StructFilteredCtx(context.Background(), s, fn)
}







func (v *Validate) StructFilteredCtx(ctx context.Context, s interface{}, fn FilterFunc) (err error) {
	val := reflect.ValueOf(s)
	top := val

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct || val.Type().ConvertibleTo(timeType) {
		return &InvalidValidationError{Type: reflect.TypeOf(s)}
	}

	
	vd := v.pool.Get().(*validate)
	vd.top = top
	vd.isPartial = true
	vd.ffn = fn
	

	vd.validateStruct(ctx, top, val, val.Type(), vd.ns[0:0], vd.actualNs[0:0], nil)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}

	v.pool.Put(vd)

	return
}







func (v *Validate) StructPartial(s interface{}, fields ...string) error {
	return v.StructPartialCtx(context.Background(), s, fields...)
}








func (v *Validate) StructPartialCtx(ctx context.Context, s interface{}, fields ...string) (err error) {
	val := reflect.ValueOf(s)
	top := val

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct || val.Type().ConvertibleTo(timeType) {
		return &InvalidValidationError{Type: reflect.TypeOf(s)}
	}

	
	vd := v.pool.Get().(*validate)
	vd.top = top
	vd.isPartial = true
	vd.ffn = nil
	vd.hasExcludes = false
	vd.includeExclude = make(map[string]struct{})

	typ := val.Type()
	name := typ.Name()

	for _, k := range fields {

		flds := strings.Split(k, namespaceSeparator)
		if len(flds) > 0 {

			vd.misc = append(vd.misc[0:0], name...)
			
			if len(vd.misc) != 0 {
				vd.misc = append(vd.misc, '.')
			}

			for _, s := range flds {

				idx := strings.Index(s, leftBracket)

				if idx != -1 {
					for idx != -1 {
						vd.misc = append(vd.misc, s[:idx]...)
						vd.includeExclude[string(vd.misc)] = struct{}{}

						idx2 := strings.Index(s, rightBracket)
						idx2++
						vd.misc = append(vd.misc, s[idx:idx2]...)
						vd.includeExclude[string(vd.misc)] = struct{}{}
						s = s[idx2:]
						idx = strings.Index(s, leftBracket)
					}
				} else {

					vd.misc = append(vd.misc, s...)
					vd.includeExclude[string(vd.misc)] = struct{}{}
				}

				vd.misc = append(vd.misc, '.')
			}
		}
	}

	vd.validateStruct(ctx, top, val, typ, vd.ns[0:0], vd.actualNs[0:0], nil)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}

	v.pool.Put(vd)

	return
}







func (v *Validate) StructExcept(s interface{}, fields ...string) error {
	return v.StructExceptCtx(context.Background(), s, fields...)
}








func (v *Validate) StructExceptCtx(ctx context.Context, s interface{}, fields ...string) (err error) {
	val := reflect.ValueOf(s)
	top := val

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct || val.Type().ConvertibleTo(timeType) {
		return &InvalidValidationError{Type: reflect.TypeOf(s)}
	}

	
	vd := v.pool.Get().(*validate)
	vd.top = top
	vd.isPartial = true
	vd.ffn = nil
	vd.hasExcludes = true
	vd.includeExclude = make(map[string]struct{})

	typ := val.Type()
	name := typ.Name()

	for _, key := range fields {

		vd.misc = vd.misc[0:0]

		if len(name) > 0 {
			vd.misc = append(vd.misc, name...)
			vd.misc = append(vd.misc, '.')
		}

		vd.misc = append(vd.misc, key...)
		vd.includeExclude[string(vd.misc)] = struct{}{}
	}

	vd.validateStruct(ctx, top, val, typ, vd.ns[0:0], vd.actualNs[0:0], nil)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}

	v.pool.Put(vd)

	return
}














func (v *Validate) Var(field interface{}, tag string) error {
	return v.VarCtx(context.Background(), field, tag)
}















func (v *Validate) VarCtx(ctx context.Context, field interface{}, tag string) (err error) {
	if len(tag) == 0 || tag == skipValidationTag {
		return nil
	}

	ctag := v.fetchCacheTag(tag)

	val := reflect.ValueOf(field)
	vd := v.pool.Get().(*validate)
	vd.top = val
	vd.isPartial = false
	vd.traverseField(ctx, val, val, vd.ns[0:0], vd.actualNs[0:0], defaultCField, ctag)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}
	v.pool.Put(vd)
	return
}















func (v *Validate) VarWithValue(field interface{}, other interface{}, tag string) error {
	return v.VarWithValueCtx(context.Background(), field, other, tag)
}
















func (v *Validate) VarWithValueCtx(ctx context.Context, field interface{}, other interface{}, tag string) (err error) {
	if len(tag) == 0 || tag == skipValidationTag {
		return nil
	}
	ctag := v.fetchCacheTag(tag)
	otherVal := reflect.ValueOf(other)
	vd := v.pool.Get().(*validate)
	vd.top = otherVal
	vd.isPartial = false
	vd.traverseField(ctx, otherVal, reflect.ValueOf(field), vd.ns[0:0], vd.actualNs[0:0], defaultCField, ctag)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}
	v.pool.Put(vd)
	return
}
