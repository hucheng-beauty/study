package validation

func Limit(offset, limit *int) {

	if offset == nil || (*offset) <= 0 {
		*offset = 0
	}

	if limit == nil || (*limit) <= 0 {
		*limit = 10
	}

	if (*limit) > 100 {
		*limit = 100
	}

	// offset:页数  limit:每页数目
	if 0 != (*offset) && (*limit) != 0 {
		*offset = (*offset) * (*limit)
	}

	return
}

func Unique(src []int) (des []int) {

	if len(src) == 0 {
		return
	}

	m := make(map[int]bool)
	for _, v := range src {
		if _, ok := m[v]; !ok {
			m[v] = true
			des = append(des, v)
		}
	}

	return des
}

func UniqueStr(src []string) (des []string) {

	if len(src) == 0 {
		return
	}

	keys := make(map[string]bool)
	for _, v := range src {
		if _, ok := keys[v]; !ok {
			keys[v] = true
			des = append(des, v)
		}
	}

	return des
}
