package context

type Linked[T comparable] struct {
	K    T
	Next *Linked[T]
	h    int
}

func NewLinked[T comparable](k T) *Linked[T] {
	return &Linked[T]{K: k, h: 1}
}

func (l *Linked[T]) Add(k T) *Linked[T] {
	if l.Empty() {
		l.K = k
		l.h++
		return l
	}

	l.h++
	var next = l
	for next.Next != nil {
		next = next.Next
	}

	next.Next = NewLinked(k)
	return l
}

func (l *Linked[T]) Rem(k T) *Linked[T] {
	if l.Empty() {
		return l
	}

	if l.K == k {
		if l.Next != nil {
			l.K = l.Next.K
			l.Next = l.Next.Next
		}
		l.h--
		return l
	}

	var next = l.Next
	var prev = l
	for next != nil {
		if next.K == k {
			prev.Next = next.Next
			l.h--
			return l
		}

		prev = next
		next = next.Next
	}

	return l
}

func (l *Linked[T]) Empty() bool {
	return l.h == 0
}

func (l *Linked[T]) Len() int {
	return l.h
}

func (l *Linked[T]) Range(f func(k T)) *Linked[T] {
	if l.Empty() {
		return l
	}

	var next = l
	for next != nil {
		f(next.K)
		next = next.Next
	}

	return l
}

func (l *Linked[T]) Values() []T {
	values := make([]T, l.h)
	index := 0
	l.Range(func(k T) {
		values[index] = k
		index++
	})
	return values
}
