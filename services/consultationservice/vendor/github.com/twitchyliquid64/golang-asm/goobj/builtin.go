



package goobj







func NBuiltin() int {
	return len(builtins)
}



func BuiltinName(i int) (string, int) {
	return builtins[i].name, builtins[i].abi
}



func BuiltinIdx(name string, abi int) int {
	i, ok := builtinMap[name]
	if !ok {
		return -1
	}
	if builtins[i].abi != abi {
		return -1
	}
	return i
}

//go:generate go run mkbuiltin.go

var builtinMap map[string]int

func init() {
	builtinMap = make(map[string]int, len(builtins))
	for i, b := range builtins {
		builtinMap[b.name] = i
	}
}
