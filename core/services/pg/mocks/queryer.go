// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	sql "database/sql"

	sqlx "github.com/smartcontractkit/sqlx"
)

// Queryer is an autogenerated mock type for the Queryer type
type Queryer struct {
	mock.Mock
}

// BindNamed provides a mock function with given fields: _a0, _a1
func (_m *Queryer) BindNamed(_a0 string, _a1 interface{}) (string, []interface{}, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	var r1 []interface{}
	var r2 error
	if rf, ok := ret.Get(0).(func(string, interface{}) (string, []interface{}, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(string, interface{}) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, interface{}) []interface{}); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]interface{})
		}
	}

	if rf, ok := ret.Get(2).(func(string, interface{}) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// DriverName provides a mock function with given fields:
func (_m *Queryer) DriverName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Exec provides a mock function with given fields: query, args
func (_m *Queryer) Exec(query string, args ...interface{}) (sql.Result, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 sql.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) (sql.Result, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) sql.Result); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sql.Result)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ExecContext provides a mock function with given fields: ctx, query, args
func (_m *Queryer) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 sql.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (sql.Result, error)); ok {
		return rf(ctx, query, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) sql.Result); ok {
		r0 = rf(ctx, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sql.Result)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: dest, query, args
func (_m *Queryer) Get(dest interface{}, query string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, dest, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, string, ...interface{}) error); ok {
		r0 = rf(dest, query, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetContext provides a mock function with given fields: ctx, dest, query, args
func (_m *Queryer) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, ctx, dest, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, string, ...interface{}) error); ok {
		r0 = rf(ctx, dest, query, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NamedExec provides a mock function with given fields: query, arg
func (_m *Queryer) NamedExec(query string, arg interface{}) (sql.Result, error) {
	ret := _m.Called(query, arg)

	var r0 sql.Result
	var r1 error
	if rf, ok := ret.Get(0).(func(string, interface{}) (sql.Result, error)); ok {
		return rf(query, arg)
	}
	if rf, ok := ret.Get(0).(func(string, interface{}) sql.Result); ok {
		r0 = rf(query, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sql.Result)
		}
	}

	if rf, ok := ret.Get(1).(func(string, interface{}) error); ok {
		r1 = rf(query, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NamedQuery provides a mock function with given fields: query, arg
func (_m *Queryer) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	ret := _m.Called(query, arg)

	var r0 *sqlx.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(string, interface{}) (*sqlx.Rows, error)); ok {
		return rf(query, arg)
	}
	if rf, ok := ret.Get(0).(func(string, interface{}) *sqlx.Rows); ok {
		r0 = rf(query, arg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(string, interface{}) error); ok {
		r1 = rf(query, arg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Prepare provides a mock function with given fields: query
func (_m *Queryer) Prepare(query string) (*sql.Stmt, error) {
	ret := _m.Called(query)

	var r0 *sql.Stmt
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*sql.Stmt, error)); ok {
		return rf(query)
	}
	if rf, ok := ret.Get(0).(func(string) *sql.Stmt); ok {
		r0 = rf(query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Stmt)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareContext provides a mock function with given fields: ctx, query
func (_m *Queryer) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	ret := _m.Called(ctx, query)

	var r0 *sql.Stmt
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*sql.Stmt, error)); ok {
		return rf(ctx, query)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *sql.Stmt); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Stmt)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareNamed provides a mock function with given fields: query
func (_m *Queryer) PrepareNamed(query string) (*sqlx.NamedStmt, error) {
	ret := _m.Called(query)

	var r0 *sqlx.NamedStmt
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*sqlx.NamedStmt, error)); ok {
		return rf(query)
	}
	if rf, ok := ret.Get(0).(func(string) *sqlx.NamedStmt); ok {
		r0 = rf(query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.NamedStmt)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: query, args
func (_m *Queryer) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sql.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) (*sql.Rows, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *sql.Rows); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryContext provides a mock function with given fields: ctx, query, args
func (_m *Queryer) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sql.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (*sql.Rows, error)); ok {
		return rf(ctx, query, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) *sql.Rows); ok {
		r0 = rf(ctx, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryRow provides a mock function with given fields: query, args
func (_m *Queryer) QueryRow(query string, args ...interface{}) *sql.Row {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sql.Row
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *sql.Row); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Row)
		}
	}

	return r0
}

// QueryRowContext provides a mock function with given fields: ctx, query, args
func (_m *Queryer) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sql.Row
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) *sql.Row); ok {
		r0 = rf(ctx, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Row)
		}
	}

	return r0
}

// QueryRowx provides a mock function with given fields: query, args
func (_m *Queryer) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sqlx.Row
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *sqlx.Row); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Row)
		}
	}

	return r0
}

// QueryRowxContext provides a mock function with given fields: ctx, query, args
func (_m *Queryer) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sqlx.Row
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) *sqlx.Row); ok {
		r0 = rf(ctx, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Row)
		}
	}

	return r0
}

// Queryx provides a mock function with given fields: query, args
func (_m *Queryer) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sqlx.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) (*sqlx.Rows, error)); ok {
		return rf(query, args...)
	}
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *sqlx.Rows); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryxContext provides a mock function with given fields: ctx, query, args
func (_m *Queryer) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sqlx.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (*sqlx.Rows, error)); ok {
		return rf(ctx, query, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) *sqlx.Rows); ok {
		r0 = rf(ctx, query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sqlx.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Rebind provides a mock function with given fields: _a0
func (_m *Queryer) Rebind(_a0 string) string {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Select provides a mock function with given fields: dest, query, args
func (_m *Queryer) Select(dest interface{}, query string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, dest, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, string, ...interface{}) error); ok {
		r0 = rf(dest, query, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SelectContext provides a mock function with given fields: ctx, dest, query, args
func (_m *Queryer) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, ctx, dest, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, string, ...interface{}) error); ok {
		r0 = rf(ctx, dest, query, args...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewQueryer interface {
	mock.TestingT
	Cleanup(func())
}

// NewQueryer creates a new instance of Queryer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewQueryer(t mockConstructorTestingTNewQueryer) *Queryer {
	mock := &Queryer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
