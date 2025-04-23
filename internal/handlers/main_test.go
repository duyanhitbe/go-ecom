package handlers

type testcase[T any] struct {
	name         string
	input        T
	mockBehavior func()
	expectedCode int
}
