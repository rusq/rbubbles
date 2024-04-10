package display

type Stack[T any] []T

func (s *Stack[T]) Push(v T) {
	*s = append(*s, v)
}

func (s *Stack[T]) Pop() T {
	var empty T
	if len(*s) == 0 {
		return empty
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s Stack[T]) Len() int {
	return len(s)
}

func (s Stack[T]) Peek() T {
	return s[len(s)-1]
}
