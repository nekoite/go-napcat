package utils

type Void struct{}

type Set[T comparable] map[T]Void

func NewSet[T comparable]() Set[T] {
	return make(map[T]Void)
}

func (s Set[T]) Add(item T) {
	s[item] = Void{}
}

func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) Remove(item T) {
	delete(s, item)
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Clear() {
	for k := range s {
		delete(s, k)
	}
}

func (s Set[T]) ToSlice() []T {
	var res []T
	for k := range s {
		res = append(res, k)
	}
	return res
}
