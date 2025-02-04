package util

func MapKeys[TKey comparable, TValue any](m map[TKey]TValue) []TKey {
	keys := make([]TKey, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func MapValues[TKey comparable, TValue interface{}](m map[TKey]TValue) []TValue {
	values := make([]TValue, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
