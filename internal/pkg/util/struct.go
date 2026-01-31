package util

func Ptr[T any](t T) *T {
	return &t
}

func Value[T any](t *T) T {
	return *t
}
