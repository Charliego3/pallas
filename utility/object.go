package utility

// Nils all values is nil return true otherwise return false
func Nils(objs ...any) bool {
	for _, obj := range objs {
		if obj != nil {
			return false
		}
	}
	return true
}

// DObj returns source if not nil else return default value
func DObj[T any](source, d T) T {
	if (any)(source) == nil {
		return d
	}
	return source
}
