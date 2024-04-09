package tools

func Map[T, V any](ts []T, fn func(T) (V, error)) ([]V, error) {
	result := make([]V, len(ts))
	for i, t := range ts {
		res, err := fn(t)
		if err != nil {
			return make([]V, 0), err
		}
		result[i] = res
	}
	return result, nil
}
