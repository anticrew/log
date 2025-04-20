package maps

func Invert[K comparable, V comparable, I map[K]V, O map[V]K](i I) O {
	if i == nil {
		return nil
	}

	result := make(O, len(i))

	for k, v := range i {
		result[v] = k
	}

	return result
}
