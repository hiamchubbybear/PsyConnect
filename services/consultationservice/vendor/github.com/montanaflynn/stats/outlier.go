package stats


type Outliers struct {
	Mild    Float64Data
	Extreme Float64Data
}


func QuartileOutliers(input Float64Data) (Outliers, error) {
	if input.Len() == 0 {
		return Outliers{}, EmptyInputErr
	}

	
	copy := sortedCopy(input)

	
	qs, _ := Quartile(copy)
	iqr, _ := InterQuartileRange(copy)

	
	lif := qs.Q1 - (1.5 * iqr)
	uif := qs.Q3 + (1.5 * iqr)
	lof := qs.Q1 - (3 * iqr)
	uof := qs.Q3 + (3 * iqr)

	
	
	
	var mild Float64Data
	var extreme Float64Data
	for _, v := range copy {

		if v < lof || v > uof {
			extreme = append(extreme, v)
		} else if v < lif || v > uif {
			mild = append(mild, v)
		}
	}

	
	return Outliers{mild, extreme}, nil
}
