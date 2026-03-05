





package mongo

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)






type WriteModel interface {
	writeModel()
}


type InsertOneModel struct {
	Document interface{}
}


func NewInsertOneModel() *InsertOneModel {
	return &InsertOneModel{}
}




func (iom *InsertOneModel) SetDocument(doc interface{}) *InsertOneModel {
	iom.Document = doc
	return iom
}

func (*InsertOneModel) writeModel() {}


type DeleteOneModel struct {
	Filter    interface{}
	Collation *options.Collation
	Hint      interface{}
}


func NewDeleteOneModel() *DeleteOneModel {
	return &DeleteOneModel{}
}




func (dom *DeleteOneModel) SetFilter(filter interface{}) *DeleteOneModel {
	dom.Filter = filter
	return dom
}



func (dom *DeleteOneModel) SetCollation(collation *options.Collation) *DeleteOneModel {
	dom.Collation = collation
	return dom
}







func (dom *DeleteOneModel) SetHint(hint interface{}) *DeleteOneModel {
	dom.Hint = hint
	return dom
}

func (*DeleteOneModel) writeModel() {}


type DeleteManyModel struct {
	Filter    interface{}
	Collation *options.Collation
	Hint      interface{}
}


func NewDeleteManyModel() *DeleteManyModel {
	return &DeleteManyModel{}
}



func (dmm *DeleteManyModel) SetFilter(filter interface{}) *DeleteManyModel {
	dmm.Filter = filter
	return dmm
}



func (dmm *DeleteManyModel) SetCollation(collation *options.Collation) *DeleteManyModel {
	dmm.Collation = collation
	return dmm
}







func (dmm *DeleteManyModel) SetHint(hint interface{}) *DeleteManyModel {
	dmm.Hint = hint
	return dmm
}

func (*DeleteManyModel) writeModel() {}


type ReplaceOneModel struct {
	Collation   *options.Collation
	Upsert      *bool
	Filter      interface{}
	Replacement interface{}
	Hint        interface{}
}


func NewReplaceOneModel() *ReplaceOneModel {
	return &ReplaceOneModel{}
}







func (rom *ReplaceOneModel) SetHint(hint interface{}) *ReplaceOneModel {
	rom.Hint = hint
	return rom
}




func (rom *ReplaceOneModel) SetFilter(filter interface{}) *ReplaceOneModel {
	rom.Filter = filter
	return rom
}



func (rom *ReplaceOneModel) SetReplacement(rep interface{}) *ReplaceOneModel {
	rom.Replacement = rep
	return rom
}



func (rom *ReplaceOneModel) SetCollation(collation *options.Collation) *ReplaceOneModel {
	rom.Collation = collation
	return rom
}




func (rom *ReplaceOneModel) SetUpsert(upsert bool) *ReplaceOneModel {
	rom.Upsert = &upsert
	return rom
}

func (*ReplaceOneModel) writeModel() {}


type UpdateOneModel struct {
	Collation    *options.Collation
	Upsert       *bool
	Filter       interface{}
	Update       interface{}
	ArrayFilters *options.ArrayFilters
	Hint         interface{}
}


func NewUpdateOneModel() *UpdateOneModel {
	return &UpdateOneModel{}
}







func (uom *UpdateOneModel) SetHint(hint interface{}) *UpdateOneModel {
	uom.Hint = hint
	return uom
}




func (uom *UpdateOneModel) SetFilter(filter interface{}) *UpdateOneModel {
	uom.Filter = filter
	return uom
}



func (uom *UpdateOneModel) SetUpdate(update interface{}) *UpdateOneModel {
	uom.Update = update
	return uom
}



func (uom *UpdateOneModel) SetArrayFilters(filters options.ArrayFilters) *UpdateOneModel {
	uom.ArrayFilters = &filters
	return uom
}



func (uom *UpdateOneModel) SetCollation(collation *options.Collation) *UpdateOneModel {
	uom.Collation = collation
	return uom
}




func (uom *UpdateOneModel) SetUpsert(upsert bool) *UpdateOneModel {
	uom.Upsert = &upsert
	return uom
}

func (*UpdateOneModel) writeModel() {}


type UpdateManyModel struct {
	Collation    *options.Collation
	Upsert       *bool
	Filter       interface{}
	Update       interface{}
	ArrayFilters *options.ArrayFilters
	Hint         interface{}
}


func NewUpdateManyModel() *UpdateManyModel {
	return &UpdateManyModel{}
}







func (umm *UpdateManyModel) SetHint(hint interface{}) *UpdateManyModel {
	umm.Hint = hint
	return umm
}



func (umm *UpdateManyModel) SetFilter(filter interface{}) *UpdateManyModel {
	umm.Filter = filter
	return umm
}



func (umm *UpdateManyModel) SetUpdate(update interface{}) *UpdateManyModel {
	umm.Update = update
	return umm
}



func (umm *UpdateManyModel) SetArrayFilters(filters options.ArrayFilters) *UpdateManyModel {
	umm.ArrayFilters = &filters
	return umm
}



func (umm *UpdateManyModel) SetCollation(collation *options.Collation) *UpdateManyModel {
	umm.Collation = collation
	return umm
}




func (umm *UpdateManyModel) SetUpsert(upsert bool) *UpdateManyModel {
	umm.Upsert = &upsert
	return umm
}

func (*UpdateManyModel) writeModel() {}
