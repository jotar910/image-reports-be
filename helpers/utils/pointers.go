package utils

func Pointer[T any](a T) *T {
	return &a
}

func NoPointerE[TValue any, TError any](a *TValue, err TError) (TValue, TError) {
	if a != nil {
		return *a, err
	}
	return *new(TValue), err
}
