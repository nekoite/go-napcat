package event

import "github.com/nekoite/go-napcat/errors"

func GetAs[T any](e any) *T {
	t, ok := e.(*T)
	if !ok {
		return nil
	}
	return t
}

func GetAsUnsafe[T any](e any) *T {
	return e.(*T)
}

func GetAsOrError[T any](e any) (*T, error) {
	t, ok := e.(*T)
	if !ok {
		return nil, errors.ErrTypeAssertion
	}
	return t, nil
}
