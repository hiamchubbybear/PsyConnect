



package language



func BaseLanguages() []Language {
	base := make([]Language, 0, NumLanguages)
	for i := 0; i < langNoIndexOffset; i++ {
		
		if i != nonCanonicalUnd {
			base = append(base, Language(i))
		}
	}
	i := langNoIndexOffset
	for _, v := range langNoIndex {
		for k := 0; k < 8; k++ {
			if v&1 == 1 {
				base = append(base, Language(i))
			}
			v >>= 1
			i++
		}
	}
	return base
}
