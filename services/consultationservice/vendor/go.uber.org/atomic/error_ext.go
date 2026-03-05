



















package atomic




//go:generate bin/gen-atomicwrapper -name=Error -type=error -wrapped=Value -pack=packError -unpack=unpackError -file=error.go

type packedError struct{ Value error }

func packError(v error) interface{} {
	return packedError{v}
}

func unpackError(v interface{}) error {
	if err, ok := v.(packedError); ok {
		return err.Value
	}
	return nil
}
