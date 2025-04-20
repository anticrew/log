package maps

func Clear[K comparable, V comparable, M map[K]V](m M) M {
	if m == nil {
		return nil
	}

	for k := range m {
		delete(m, k)
	}

	return m
}
