package handlers

type testcase[T any] struct {
	name         string
	path         string
	input        T
	mockBehavior func()
	expectedCode int
}
