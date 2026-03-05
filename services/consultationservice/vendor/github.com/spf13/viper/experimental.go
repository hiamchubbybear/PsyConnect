package viper


func ExperimentalBindStruct() Option {
	return optionFunc(func(v *Viper) {
		v.experimentalBindStruct = true
	})
}
