package queue

// Queue simple ring queue with a growing capacity
type Queue[T any] struct {
	buf   []T
	front int
	size  int
}

// Size returns the number of items in the queue
func (q *Queue[T]) Size() int {
	return q.size
}

// Front returns the first item from the queue. If the queue size is zero, Front panics
func (q *Queue[T]) Front() T {
	if q.size == 0 {
		panic("Queue.Front: queue is empty")
	}

	return q.buf[q.front]
}

// Pop returns and removes the first item from the queue. If the queue size is zero, Pop panics
func (q *Queue[T]) Pop() T {
	v := q.Front()
	q.size--

	q.front++
	if q.front == len(q.buf) {
		q.front = 0
	}

	return v
}

// Push appends an item to the back of the queue. Push grows queue capacity if necessary
func (q *Queue[T]) Push(v T) {
	if q.size == len(q.buf) {
		q.grow(1)
	}

	back := q.front + q.size
	if back >= len(q.buf) {
		back -= len(q.buf)
	}

	q.buf[back] = v
	q.size++
}

func (q *Queue[T]) grow(n int) {
	newBuf := make([]T, 2*len(q.buf)+n)

	copy(newBuf, q.buf[q.front:])
	copy(newBuf[len(q.buf)-q.front:], q.buf)

	q.buf = newBuf
	q.front = 0
}

// Grow grows queue capacity, if necessary, to guarantee space for another n items.
// After Grow(n), at least n items can be pushed to queue without another allocation.
// If n is negative, Grow panics.
func (q *Queue[T]) Grow(n int) {
	if n < 0 {
		panic("Queue.Grow: negative count")
	} else if len(q.buf)-q.size < n {
		q.grow(n)
	}
}

// Items returns slice of queue items. Never returns nil.
func (q *Queue[T]) Items() []T {
	sl := make([]T, q.size)
	n := min(q.size, len(q.buf)-q.front)
	copy(sl, q.buf[q.front:n])
	copy(sl[n:], q.buf)
	return sl
}

func (q *Queue[T]) Clear() {
	q.front = 0
	q.size = 0
}

// Deque two-sided ring queue with a growing capacity. This Queue extension to get last
// and add first item to queue. Deque provides all public methods of Queue
type Deque[T any] struct {
	Queue[T]
}

// Back returns the last item from the queue. If the queue size is zero, Back panics
func (q *Deque[T]) Back() T {
	if q.size == 0 {
		panic("Queue.Back: queue is empty")
	}

	back := q.front + q.size - 1
	if back >= len(q.buf) {
		back -= len(q.buf)
	}

	return q.buf[back]
}

// PopFront is the same as Queue.Pop
func (q *Deque[T]) PopFront() T {
	return q.Pop()
}

// PopBack returns and removes the last item from the queue. If the queue size is zero,
// PopBack panics
func (q *Deque[T]) PopBack() T {
	v := q.Back()
	q.size--
	return v
}

// PushFront adds an item to the front of the queue. PushFront grows queue capacity if necessary
func (q *Deque[T]) PushFront(v T) {
	if q.size == len(q.buf) {
		q.grow(1)
	}

	if q.front > 0 {
		q.front--
	} else {
		q.front = len(q.buf) - 1
	}

	q.buf[q.front] = v
	q.size++
}

// PushBack is the same as Queue.Push
func (q *Deque[T]) PushBack(v T) {
	q.Push(v)
}

type Int interface {
	~int | ~int64 | ~int32 | ~int16 | ~int8
}

type Uint interface {
	~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

type Float interface {
	~float64 | ~float32
}

type Ordered interface {
	Int | Uint | Float | ~string
}

func min[T Ordered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

func max[T Ordered](a, b T) T {
	if a < b {
		return b
	}
	return a
}
