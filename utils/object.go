package utils

// Nils all values is nil return true otherwise return false
func Nils(objs ...any) bool {
	for _, obj := range objs {
		if obj != nil {
			return false
		}
	}
	return true
}
