//go:build go1.21 && !go1.24
// +build go1.21,!go1.24



package resolver

import (
    _ `encoding/json`
    `reflect`
    _ `unsafe`
)

type StdField struct {
    name        string
    nameBytes   []byte
    nameNonEsc  string
    nameEscHTML string
    tag         bool
    index       []int
    typ         reflect.Type
    omitEmpty   bool
    quoted      bool
    encoder     func()
}

type StdStructFields struct {
    list      []StdField
    nameIndex map[string]*StdField
    byFoldedName map[string]*StdField
}

//go:noescape
//go:linkname typeFields encoding/json.typeFields
func typeFields(_ reflect.Type) StdStructFields

func handleOmitZero(f StdField, fv *FieldMeta) {}
