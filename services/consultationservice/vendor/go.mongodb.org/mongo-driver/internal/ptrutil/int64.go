





package ptrutil








func CompareInt64(ptr1, ptr2 *int64) int {
	if ptr1 == ptr2 {
		
		return 0
	}

	if ptr1 == nil && ptr2 != nil {
		return -2
	}

	if ptr1 != nil && ptr2 == nil {
		return 2
	}

	if *ptr1 > *ptr2 {
		return 1
	}

	if *ptr1 < *ptr2 {
		return -1
	}

	return 0
}
