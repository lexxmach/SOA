package common

// Panics on err, else returns v
func Must[T any](v T, err error) T {
	if err != nil {
		panic("error: " + err.Error())
	}
	return v
}
