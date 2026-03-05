package validator

import "reflect"



type FieldLevel interface {

	
	Top() reflect.Value

	
	
	Parent() reflect.Value

	
	Field() reflect.Value

	
	
	FieldName() string

	
	StructFieldName() string

	
	Param() string

	
	GetTag() string

	
	
	
	ExtractType(field reflect.Value) (value reflect.Value, kind reflect.Kind, nullable bool)

	
	
	
	
	
	
	
	
	GetStructFieldOK() (reflect.Value, reflect.Kind, bool)

	
	
	
	
	GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool)

	
	
	
	
	
	
	GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool)

	
	
	GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool)
}

var _ FieldLevel = new(validate)


func (v *validate) Field() reflect.Value {
	return v.flField
}



func (v *validate) FieldName() string {
	return v.cf.altName
}


func (v *validate) GetTag() string {
	return v.ct.tag
}


func (v *validate) StructFieldName() string {
	return v.cf.name
}


func (v *validate) Param() string {
	return v.ct.param
}




func (v *validate) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	current, kind, _, found := v.getStructFieldOKInternal(v.slflParent, v.ct.param)
	return current, kind, found
}





func (v *validate) GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) {
	current, kind, _, found := v.GetStructFieldOKAdvanced2(val, namespace)
	return current, kind, found
}


func (v *validate) GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool) {
	return v.getStructFieldOKInternal(v.slflParent, v.ct.param)
}



func (v *validate) GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool) {
	return v.getStructFieldOKInternal(val, namespace)
}
