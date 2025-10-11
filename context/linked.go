package context

type Linked[T comparable] struct {
	K    T
	Next *Linked[T]
	h    int
}

func NewLinked[T comparable](k T) *Linked[T] {
	return &Linked[T]{K: k, h: 1}
}

func (l *Linked[T]) Add(k T) {
	l.h++
	if l.Next == nil {
		l.Next = NewLinked(k)
		return
	}

	var next = l.Next
	for next.Next != nil {
		next = next.Next
	}

	next.Next = NewLinked(k)
}

func (l *Linked[T]) Rem(k T) {
	if l.K == k {
		if l.Next != nil {
			l.K = l.Next.K
			l.Next = l.Next.Next
			l.h--
		}
		return
	}

	var next = l.Next
	var prev = l
	for {
		if next.K == k {
			prev.Next = next.Next
			l.h--
			return
		}

		prev = next
		next = next.Next
	}
}

func (l *Linked[T]) Len() int {
	return l.h
}

func (l *Linked[T]) Range(f func(k T)) {
	var next = l
	for next != nil {
		f(next.K)
		next = next.Next
	}
}

func (l *Linked[T]) Values() []T {
	values := make([]T, l.h)
	var next = l
	index := 0
	for next != nil {
		values[index] = next.K
		index++
		next = next.Next
	}

	return values
}
