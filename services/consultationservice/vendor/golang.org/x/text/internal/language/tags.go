



package language



func MustParse(s string) Tag {
	t, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return t
}



func MustParseBase(s string) Language {
	b, err := ParseBase(s)
	if err != nil {
		panic(err)
	}
	return b
}



func MustParseScript(s string) Script {
	scr, err := ParseScript(s)
	if err != nil {
		panic(err)
	}
	return scr
}



func MustParseRegion(s string) Region {
	r, err := ParseRegion(s)
	if err != nil {
		panic(err)
	}
	return r
}


var Und Tag
