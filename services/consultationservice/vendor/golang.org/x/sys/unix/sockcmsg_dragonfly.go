



package unix


func cmsgAlignOf(salen int) int {
	salign := SizeofPtr
	if SizeofPtr == 8 && !supportsABI(_dragonflyABIChangeVersion) {
		
		
		salign = 4
	}
	return (salen + salign - 1) & ^(salign - 1)
}
