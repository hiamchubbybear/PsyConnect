


package attribute 






type Filter func(KeyValue) bool





func NewAllowKeysFilter(keys ...Key) Filter {
	if len(keys) <= 0 {
		return func(kv KeyValue) bool { return false }
	}

	allowed := make(map[Key]struct{})
	for _, k := range keys {
		allowed[k] = struct{}{}
	}
	return func(kv KeyValue) bool {
		_, ok := allowed[kv.Key]
		return ok
	}
}





func NewDenyKeysFilter(keys ...Key) Filter {
	if len(keys) <= 0 {
		return func(kv KeyValue) bool { return true }
	}

	forbid := make(map[Key]struct{})
	for _, k := range keys {
		forbid[k] = struct{}{}
	}
	return func(kv KeyValue) bool {
		_, ok := forbid[kv.Key]
		return !ok
	}
}
