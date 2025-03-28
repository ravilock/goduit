// Code generated by mockery v2.53.2. DO NOT EDIT.

package services

import (
	context "context"

	models "github.com/ravilock/goduit/internal/profileManager/models"
	mock "github.com/stretchr/testify/mock"
)

// mockProfileUpdater is an autogenerated mock type for the profileUpdater type
type mockProfileUpdater struct {
	mock.Mock
}

type mockProfileUpdater_Expecter struct {
	mock *mock.Mock
}

func (_m *mockProfileUpdater) EXPECT() *mockProfileUpdater_Expecter {
	return &mockProfileUpdater_Expecter{mock: &_m.Mock}
}

// UpdateProfile provides a mock function with given fields: ctx, subjectEmail, clientUsername, user
func (_m *mockProfileUpdater) UpdateProfile(ctx context.Context, subjectEmail string, clientUsername string, user *models.User) error {
	ret := _m.Called(ctx, subjectEmail, clientUsername, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProfile")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *models.User) error); ok {
		r0 = rf(ctx, subjectEmail, clientUsername, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockProfileUpdater_UpdateProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateProfile'
type mockProfileUpdater_UpdateProfile_Call struct {
	*mock.Call
}

// UpdateProfile is a helper method to define mock.On call
//   - ctx context.Context
//   - subjectEmail string
//   - clientUsername string
//   - user *models.User
func (_e *mockProfileUpdater_Expecter) UpdateProfile(ctx interface{}, subjectEmail interface{}, clientUsername interface{}, user interface{}) *mockProfileUpdater_UpdateProfile_Call {
	return &mockProfileUpdater_UpdateProfile_Call{Call: _e.mock.On("UpdateProfile", ctx, subjectEmail, clientUsername, user)}
}

func (_c *mockProfileUpdater_UpdateProfile_Call) Run(run func(ctx context.Context, subjectEmail string, clientUsername string, user *models.User)) *mockProfileUpdater_UpdateProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(*models.User))
	})
	return _c
}

func (_c *mockProfileUpdater_UpdateProfile_Call) Return(_a0 error) *mockProfileUpdater_UpdateProfile_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockProfileUpdater_UpdateProfile_Call) RunAndReturn(run func(context.Context, string, string, *models.User) error) *mockProfileUpdater_UpdateProfile_Call {
	_c.Call.Return(run)
	return _c
}

// newMockProfileUpdater creates a new instance of mockProfileUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockProfileUpdater(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockProfileUpdater {
	mock := &mockProfileUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
