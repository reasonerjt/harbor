// Code generated by mockery v2.51.0. DO NOT EDIT.

package dao

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/goharbor/harbor/src/pkg/systemartifact/model"

	q "github.com/goharbor/harbor/src/lib/q"
)

// DAO is an autogenerated mock type for the DAO type
type DAO struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, systemArtifact
func (_m *DAO) Create(ctx context.Context, systemArtifact *model.SystemArtifact) (int64, error) {
	ret := _m.Called(ctx, systemArtifact)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.SystemArtifact) (int64, error)); ok {
		return rf(ctx, systemArtifact)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.SystemArtifact) int64); ok {
		r0 = rf(ctx, systemArtifact)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.SystemArtifact) error); ok {
		r1 = rf(ctx, systemArtifact)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, vendor, repository, digest
func (_m *DAO) Delete(ctx context.Context, vendor string, repository string, digest string) error {
	ret := _m.Called(ctx, vendor, repository, digest)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, vendor, repository, digest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, vendor, repository, digest
func (_m *DAO) Get(ctx context.Context, vendor string, repository string, digest string) (*model.SystemArtifact, error) {
	ret := _m.Called(ctx, vendor, repository, digest)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *model.SystemArtifact
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) (*model.SystemArtifact, error)); ok {
		return rf(ctx, vendor, repository, digest)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *model.SystemArtifact); ok {
		r0 = rf(ctx, vendor, repository, digest)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SystemArtifact)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, vendor, repository, digest)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, query
func (_m *DAO) List(ctx context.Context, query *q.Query) ([]*model.SystemArtifact, error) {
	ret := _m.Called(ctx, query)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*model.SystemArtifact
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *q.Query) ([]*model.SystemArtifact, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *q.Query) []*model.SystemArtifact); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.SystemArtifact)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *q.Query) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Size provides a mock function with given fields: ctx
func (_m *DAO) Size(ctx context.Context) (int64, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Size")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (int64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) int64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDAO creates a new instance of DAO. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDAO(t interface {
	mock.TestingT
	Cleanup(func())
}) *DAO {
	mock := &DAO{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
