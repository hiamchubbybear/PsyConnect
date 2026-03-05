





package options


type DataKeyOptions struct {
	MasterKey   interface{}
	KeyAltNames []string

	
	
	KeyMaterial []byte
}


func DataKey() *DataKeyOptions {
	return &DataKeyOptions{}
}








































func (dk *DataKeyOptions) SetMasterKey(masterKey interface{}) *DataKeyOptions {
	dk.MasterKey = masterKey
	return dk
}



func (dk *DataKeyOptions) SetKeyAltNames(keyAltNames []string) *DataKeyOptions {
	dk.KeyAltNames = keyAltNames
	return dk
}


func (dk *DataKeyOptions) SetKeyMaterial(keyMaterial []byte) *DataKeyOptions {
	dk.KeyMaterial = keyMaterial
	return dk
}





func MergeDataKeyOptions(opts ...*DataKeyOptions) *DataKeyOptions {
	dko := DataKey()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.MasterKey != nil {
			dko.MasterKey = opt.MasterKey
		}
		if opt.KeyAltNames != nil {
			dko.KeyAltNames = opt.KeyAltNames
		}
		if opt.KeyMaterial != nil {
			dko.KeyMaterial = opt.KeyMaterial
		}
	}

	return dko
}
