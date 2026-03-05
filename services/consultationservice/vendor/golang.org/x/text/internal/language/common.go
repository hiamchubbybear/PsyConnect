

package language




type AliasType int8

const (
	Deprecated AliasType = iota
	Macro
	Legacy

	AliasTypeUnknown AliasType = -1
)
