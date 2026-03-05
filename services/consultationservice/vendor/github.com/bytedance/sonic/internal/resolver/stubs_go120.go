//go:build !go1.21
// +build !go1.21



package resolver

import (
    _ `encoding/json`
    `reflect`
    _ `unsafe`
)

type StdField struct {
    name        string
    nameBytes   []byte
    equalFold   func()
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
    nameIndex map[string]int
}

//go:noescape
//go:linkname typeFields encoding/json.typeFields
func typeFields(_ reflect.Type) StdStructFields

func handleOmitZero(f StdField, fv *FieldMeta) {}
