package validator


type Option func(*Validate)







func WithRequiredStructEnabled() Option {
	return func(v *Validate) {
		v.requiredStructEnabled = true
	}
}





func WithPrivateFieldValidation() Option {
	return func(v *Validate) {
		v.privateFieldValidation = true
	}
}
