package events

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/walletera/werrors"
)

// MockHandler is a mock implementation of the Handler interface
type MockHandler struct {
	mock.Mock
}

type MockHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *MockHandler) EXPECT() *MockHandler_Expecter {
	return &MockHandler_Expecter{mock: &_m.Mock}
}

// HandlePaymentCreated provides a mock function with given fields: ctx, paymentCreatedEvent
func (_m *MockHandler) HandlePaymentCreated(ctx context.Context, paymentCreatedEvent PaymentCreated) werrors.WError {
	ret := _m.Called(ctx, paymentCreatedEvent)

	if len(ret) == 0 {
		panic("no return value specified for HandlePaymentCreated")
	}

	var r0 werrors.WError
	if rf, ok := ret.Get(0).(func(context.Context, PaymentCreated) werrors.WError); ok {
		r0 = rf(ctx, paymentCreatedEvent)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(werrors.WError)
		}
	}

	return r0
}

// MockHandler_HandlePaymentCreated_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HandlePaymentCreated'
type MockHandler_HandlePaymentCreated_Call struct {
	*mock.Call
}

// HandlePaymentCreated is a helper method to define mock.On call
//   - ctx context.Context
//   - paymentCreatedEvent PaymentCreated
func (_e *MockHandler_Expecter) HandlePaymentCreated(ctx interface{}, paymentCreatedEvent interface{}) *MockHandler_HandlePaymentCreated_Call {
	return &MockHandler_HandlePaymentCreated_Call{Call: _e.mock.On("HandlePaymentCreated", ctx, paymentCreatedEvent)}
}

func (_c *MockHandler_HandlePaymentCreated_Call) Run(run func(ctx context.Context, paymentCreatedEvent PaymentCreated)) *MockHandler_HandlePaymentCreated_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(PaymentCreated))
	})
	return _c
}

func (_c *MockHandler_HandlePaymentCreated_Call) Return(_a0 werrors.WError) *MockHandler_HandlePaymentCreated_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHandler_HandlePaymentCreated_Call) RunAndReturn(run func(context.Context, PaymentCreated) werrors.WError) *MockHandler_HandlePaymentCreated_Call {
	_c.Call.Return(run)
	return _c
}

// HandlePaymentUpdated provides a mock function with given fields: ctx, paymentUpdatedEvent
func (_m *MockHandler) HandlePaymentUpdated(ctx context.Context, paymentUpdatedEvent PaymentUpdated) werrors.WError {
	ret := _m.Called(ctx, paymentUpdatedEvent)

	if len(ret) == 0 {
		panic("no return value specified for HandlePaymentUpdated")
	}

	var r0 werrors.WError
	if rf, ok := ret.Get(0).(func(context.Context, PaymentUpdated) werrors.WError); ok {
		r0 = rf(ctx, paymentUpdatedEvent)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(werrors.WError)
		}
	}

	return r0
}

// MockHandler_HandlePaymentUpdated_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HandlePaymentUpdated'
type MockHandler_HandlePaymentUpdated_Call struct {
	*mock.Call
}

// HandlePaymentUpdated is a helper method to define mock.On call
//   - ctx context.Context
//   - paymentUpdatedEvent PaymentUpdated
func (_e *MockHandler_Expecter) HandlePaymentUpdated(ctx interface{}, paymentUpdatedEvent interface{}) *MockHandler_HandlePaymentUpdated_Call {
	return &MockHandler_HandlePaymentUpdated_Call{Call: _e.mock.On("HandlePaymentUpdated", ctx, paymentUpdatedEvent)}
}

func (_c *MockHandler_HandlePaymentUpdated_Call) Run(run func(ctx context.Context, paymentUpdatedEvent PaymentUpdated)) *MockHandler_HandlePaymentUpdated_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(PaymentUpdated))
	})
	return _c
}

func (_c *MockHandler_HandlePaymentUpdated_Call) Return(_a0 werrors.WError) *MockHandler_HandlePaymentUpdated_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockHandler_HandlePaymentUpdated_Call) RunAndReturn(run func(context.Context, PaymentUpdated) werrors.WError) *MockHandler_HandlePaymentUpdated_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockHandler creates a new instance of MockHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockHandler {
	mock := &MockHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}