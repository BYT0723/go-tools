package ds

// Stack is a generic interface for Last-In-First-Out (LIFO) data structures.
// It provides basic stack operations like push, pop, and peek.
//
// Type parameters:
//   - T: The element type stored in the stack
type Stack[T any] interface {
	// Push adds an element to the top of the stack.
	Push(T)
	// Pop removes and returns the top element from the stack.
	// Returns the element and a boolean indicating if the stack was not empty.
	Pop() (T, bool)
	// Peek returns the top element without removing it from the stack.
	// Returns the element and a boolean indicating if the stack was not empty.
	Peek() (T, bool)
	// Empty returns true if the stack contains no elements.
	Empty() bool
	// Size returns the number of elements in the stack.
	Size() int
}
