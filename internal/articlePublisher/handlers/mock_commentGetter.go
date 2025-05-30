// Code generated by mockery v2.53.2. DO NOT EDIT.

package handlers

import (
	context "context"

	models "github.com/ravilock/goduit/internal/articlePublisher/models"
	mock "github.com/stretchr/testify/mock"
)

// mockCommentGetter is an autogenerated mock type for the commentGetter type
type mockCommentGetter struct {
	mock.Mock
}

type mockCommentGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *mockCommentGetter) EXPECT() *mockCommentGetter_Expecter {
	return &mockCommentGetter_Expecter{mock: &_m.Mock}
}

// GetCommentByID provides a mock function with given fields: ctx, ID
func (_m *mockCommentGetter) GetCommentByID(ctx context.Context, ID string) (*models.Comment, error) {
	ret := _m.Called(ctx, ID)

	if len(ret) == 0 {
		panic("no return value specified for GetCommentByID")
	}

	var r0 *models.Comment
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.Comment, error)); ok {
		return rf(ctx, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Comment); ok {
		r0 = rf(ctx, ID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Comment)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockCommentGetter_GetCommentByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetCommentByID'
type mockCommentGetter_GetCommentByID_Call struct {
	*mock.Call
}

// GetCommentByID is a helper method to define mock.On call
//   - ctx context.Context
//   - ID string
func (_e *mockCommentGetter_Expecter) GetCommentByID(ctx interface{}, ID interface{}) *mockCommentGetter_GetCommentByID_Call {
	return &mockCommentGetter_GetCommentByID_Call{Call: _e.mock.On("GetCommentByID", ctx, ID)}
}

func (_c *mockCommentGetter_GetCommentByID_Call) Run(run func(ctx context.Context, ID string)) *mockCommentGetter_GetCommentByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockCommentGetter_GetCommentByID_Call) Return(_a0 *models.Comment, _a1 error) *mockCommentGetter_GetCommentByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockCommentGetter_GetCommentByID_Call) RunAndReturn(run func(context.Context, string) (*models.Comment, error)) *mockCommentGetter_GetCommentByID_Call {
	_c.Call.Return(run)
	return _c
}

// newMockCommentGetter creates a new instance of mockCommentGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockCommentGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockCommentGetter {
	mock := &mockCommentGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
