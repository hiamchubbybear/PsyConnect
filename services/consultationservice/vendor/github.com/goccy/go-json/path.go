package json

import (
	"reflect"

	"github.com/goccy/go-json/internal/decoder"
)















func CreatePath(p string) (*Path, error) {
	path, err := decoder.PathString(p).Build()
	if err != nil {
		return nil, err
	}
	return &Path{path: path}, nil
}


type Path struct {
	path *decoder.Path
}


func (p *Path) RootSelectorOnly() bool {
	return p.path.RootSelectorOnly
}


func (p *Path) UsedSingleQuotePathSelector() bool {
	return p.path.SingleQuotePathSelector
}


func (p *Path) UsedDoubleQuotePathSelector() bool {
	return p.path.DoubleQuotePathSelector
}


func (p *Path) Extract(data []byte, optFuncs ...DecodeOptionFunc) ([][]byte, error) {
	return extractFromPath(p, data, optFuncs...)
}


func (p *Path) PathString() string {
	return p.path.String()
}


func (p *Path) Unmarshal(data []byte, v interface{}, optFuncs ...DecodeOptionFunc) error {
	contents, err := extractFromPath(p, data, optFuncs...)
	if err != nil {
		return err
	}
	results := make([]interface{}, 0, len(contents))
	for _, content := range contents {
		var result interface{}
		if err := Unmarshal(content, &result); err != nil {
			return err
		}
		results = append(results, result)
	}
	if err := decoder.AssignValue(reflect.ValueOf(results), reflect.ValueOf(v)); err != nil {
		return err
	}
	return nil
}


func (p *Path) Get(src, dst interface{}) error {
	return p.path.Get(reflect.ValueOf(src), reflect.ValueOf(dst))
}
