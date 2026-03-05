





package options

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



const (
	QueryTypeEquality string = "equality"
)


type RangeOptions struct {
	Min        *bson.RawValue
	Max        *bson.RawValue
	Sparsity   *int64
	TrimFactor *int32
	Precision  *int32
}


type EncryptOptions struct {
	KeyID            *primitive.Binary
	KeyAltName       *string
	Algorithm        string
	QueryType        string
	ContentionFactor *int64
	RangeOptions     *RangeOptions
}


func Encrypt() *EncryptOptions {
	return &EncryptOptions{}
}


func (e *EncryptOptions) SetKeyID(keyID primitive.Binary) *EncryptOptions {
	e.KeyID = &keyID
	return e
}


func (e *EncryptOptions) SetKeyAltName(keyAltName string) *EncryptOptions {
	e.KeyAltName = &keyAltName
	return e
}









func (e *EncryptOptions) SetAlgorithm(algorithm string) *EncryptOptions {
	e.Algorithm = algorithm
	return e
}





func (e *EncryptOptions) SetQueryType(queryType string) *EncryptOptions {
	e.QueryType = queryType
	return e
}



func (e *EncryptOptions) SetContentionFactor(contentionFactor int64) *EncryptOptions {
	e.ContentionFactor = &contentionFactor
	return e
}


func (e *EncryptOptions) SetRangeOptions(ro RangeOptions) *EncryptOptions {
	e.RangeOptions = &ro
	return e
}


func (ro *RangeOptions) SetMin(min bson.RawValue) *RangeOptions {
	ro.Min = &min
	return ro
}


func (ro *RangeOptions) SetMax(max bson.RawValue) *RangeOptions {
	ro.Max = &max
	return ro
}


func (ro *RangeOptions) SetSparsity(sparsity int64) *RangeOptions {
	ro.Sparsity = &sparsity
	return ro
}


func (ro *RangeOptions) SetTrimFactor(trimFactor int32) *RangeOptions {
	ro.TrimFactor = &trimFactor
	return ro
}


func (ro *RangeOptions) SetPrecision(precision int32) *RangeOptions {
	ro.Precision = &precision
	return ro
}





func MergeEncryptOptions(opts ...*EncryptOptions) *EncryptOptions {
	eo := Encrypt()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.KeyID != nil {
			eo.KeyID = opt.KeyID
		}
		if opt.KeyAltName != nil {
			eo.KeyAltName = opt.KeyAltName
		}
		if opt.Algorithm != "" {
			eo.Algorithm = opt.Algorithm
		}
		if opt.QueryType != "" {
			eo.QueryType = opt.QueryType
		}
		if opt.ContentionFactor != nil {
			eo.ContentionFactor = opt.ContentionFactor
		}
		if opt.RangeOptions != nil {
			eo.RangeOptions = opt.RangeOptions
		}
	}

	return eo
}
