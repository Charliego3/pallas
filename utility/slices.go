package utility

func First[S ~[]E, E any](d E, list S) E {
	if len(list) == 0 {
		return d
	}
	return list[0]
}
