package dag_engine

func Deduplicate(in []DependAbleRunner) (out []DependAbleRunner) {
	if len(in) <= 0 {
		return in
	}

	m := make(map[DependAbleRunner]struct{})
	for _, d := range in {
		_, ok := m[d]
		if !ok {
			m[d] = struct{}{}
		}
	}

	for k, _ := range m {
		out = append(out, k)
	}
	return out
}
