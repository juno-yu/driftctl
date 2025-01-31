// Code generated by mockery v1.0.0. DO NOT EDIT.

package repository

import "github.com/stretchr/testify/mock"

// MockSQSRepository is an autogenerated mock type for the MockSQSRepository type
type MockSQSRepository struct {
	mock.Mock
}

// ListAllQueues provides a mock function with given fields:
func (_m *MockSQSRepository) ListAllQueues() ([]*string, error) {
	ret := _m.Called()

	var r0 []*string
	if rf, ok := ret.Get(0).(func() []*string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
